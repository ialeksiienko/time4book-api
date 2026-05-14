package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	httpadapter "time4book/internal/app/adapters/in/http"
	"time4book/internal/app/adapters/out/jwt"
	"time4book/internal/app/adapters/out/postgres"
	authrepo "time4book/internal/app/adapters/out/postgres/repo/auth"
	companyrepo "time4book/internal/app/adapters/out/postgres/repo/company"
	companyresourcetyperepo "time4book/internal/app/adapters/out/postgres/repo/companyresourcetype"
	bookingrepo "time4book/internal/app/adapters/out/postgres/repo/reservation"
	resourcerepo "time4book/internal/app/adapters/out/postgres/repo/resource"
	userrepo "time4book/internal/app/adapters/out/postgres/repo/user"
	"time4book/internal/app/bootstrap"
	"time4book/internal/app/config"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/usecases"
	"time4book/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Facade struct {
	commands    *usecases.Commands
	userRepo    user.UserRepo
	companyRepo company.CompanyRepo
	logger      *slog.Logger
	isProd      bool
	jwtManager  *jwt.JWTManager
}

func New() *Facade {
	cfg := config.MustLoad()

	accessDur, err := time.ParseDuration(cfg.JwtAccessTokenDuration)
	if err != nil {
		panic(fmt.Sprintf("invalid access token duration: %v", err))
	}

	refreshDur, err := time.ParseDuration(cfg.JwtRefreshTokenDuration)
	if err != nil {
		panic(fmt.Sprintf("invalid refresh token duration: %v", err))
	}

	logger := setupLogger(cfg.Env)

	dbConfig := &postgres.Config{
		User: cfg.User,
		Pass: cfg.Pass,
		Host: cfg.Host,
		Port: cfg.Port,
		Name: cfg.Name,

		Logger: logger,
	}

	if err := postgres.RunMigrations(context.Background(), dbConfig); err != nil {
		panic(fmt.Sprintf("unexpected error while trying to run database migrations: %s", err.Error()))
	}

	pool, _, err := postgres.NewDBPool(dbConfig)
	if err != nil {
		panic(fmt.Sprintf("unexpected error while trying to connect to database: %s", err.Error()))
	}

	db := postgres.New(pool)
	txManager := postgres.NewTxManager(db)
	validator := validator.New()
	jwtManager := jwt.NewManager(cfg.JwtSecret, accessDur, refreshDur)

	userRepo := userrepo.New(db)
	authRepo := authrepo.New(db)
	companyRepo := companyrepo.New(db)
	resourceRepo := resourcerepo.New(db)
	bookingRepo := bookingrepo.New(db)
	crtRepo := companyresourcetyperepo.New(db)

	commands := usecases.New(userRepo, authRepo, companyRepo, resourceRepo, crtRepo, bookingRepo, txManager, jwtManager, validator, logger)

	ctx := context.Background()
	if bootstrapErr := bootstrap.EnsureDeveloperUser(
		ctx,
		strings.TrimSpace(strings.ToLower(cfg.Env)),
		bootstrap.DeveloperBootstrapPassword,
		txManager,
		userRepo,
		authRepo,
		companyRepo,
		logger.With(slog.String("bootstrap", "developer")),
	); bootstrapErr != nil {
		logger.Warn("developer bootstrap skipped or failed",
			slog.String("error", bootstrapErr.Error()),
		)
	}

	return &Facade{
		commands:    commands,
		userRepo:    userRepo,
		companyRepo: companyRepo,
		logger:      logger,
		isProd:      cfg.Env == envProd,
		jwtManager:  jwtManager,
	}
}

func (f *Facade) GetHTTPServer() *httpadapter.Server {
	if f.isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	handler := httpadapter.NewHandler(f.commands)

	authMw := httpadapter.JWTAuth(f.jwtManager)
	companyMw := httpadapter.RequireCompanyScope(f.userRepo)
	activeCompanyMw := httpadapter.RequireActiveCompany(f.userRepo, f.companyRepo)

	router := httpadapter.NewRouter(handler, authMw, companyMw, activeCompanyMw)

	srv := httpadapter.New(f.logger, 50052, router)

	return srv
}

func setupLogger(env string) *slog.Logger {
	var (
		level     zerolog.Level
		slogLevel slog.Level
		writer    io.Writer
	)

	writer = os.Stdout
	switch strings.ToLower(env) {
	case envLocal:
		level = zerolog.DebugLevel
		slogLevel = slog.LevelDebug

		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	case envDev:
		level = zerolog.DebugLevel
		slogLevel = slog.LevelDebug

	case envProd:
		level = zerolog.InfoLevel
		slogLevel = slog.LevelInfo

	default:
		level = zerolog.InfoLevel
		slogLevel = slog.LevelInfo
	}

	zlogger := zerolog.New(writer).Level(level).With().Timestamp().Logger()
	handler := slogzerolog.Option{
		Level:  slogLevel,
		Logger: &zlogger,
	}.NewZerologHandler()

	return slog.New(handler)
}
