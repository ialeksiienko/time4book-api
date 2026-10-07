package httpadapter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv *http.Server
	log *slog.Logger
}

func New(
	log *slog.Logger,
	port int64,
	r *gin.Engine,
) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
		log: log,
	}
}

func (s *Server) Run() {
	s.log = s.log.With(slog.String("address", s.srv.Addr))

	go func() {
		s.log.Info("starting server")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("error listening: %s", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("forced shutdown: %s", err.Error()))
	}

	s.log.Info("server exited gracefully")
}
