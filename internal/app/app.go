package app

import (
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
	bookingrepo "time4book/internal/app/adapters/out/postgres/repo/reservation"
	resourcerepo "time4book/internal/app/adapters/out/postgres/repo/resource"
	userrepo "time4book/internal/app/adapters/out/postgres/repo/user"
	"time4book/internal/app/config"
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
	commands   *usecases.Commands
	logger     *slog.Logger
	isProd     bool
	jwtManager *jwt.JWTManager
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

	pool, _, err := postgres.NewDBPool(&postgres.Config{
		User: cfg.User,
		Pass: cfg.Pass,
		Host: cfg.Host,
		Port: cfg.Port,
		Name: cfg.Name,

		Logger: logger,
	})
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

	commands := usecases.New(userRepo, authRepo, companyRepo, resourceRepo, bookingRepo, txManager, jwtManager, validator, logger)

	return &Facade{
		commands:   commands,
		logger:     logger,
		isProd:     cfg.Env == envProd,
		jwtManager: jwtManager,
	}
}

func (f *Facade) GetHTTPServer() *httpadapter.Server {
	if f.isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	handler := httpadapter.NewHandler(f.commands)

	authMw := httpadapter.JWTAuth(f.jwtManager)
	companyMw := httpadapter.RequireCompany()

	router := httpadapter.NewRouter(handler, authMw, companyMw)

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
