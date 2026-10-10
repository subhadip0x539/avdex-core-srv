package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/subhadip0x539/arc-dex-core-srv/src/config"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/infra"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/router"
)

func Run(cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := infra.NewPostgres().
		Host(cfg.Postgres.Host).
		Port(cfg.Postgres.Port).
		Database(cfg.Postgres.Database).
		Username(cfg.Postgres.Username).
		Password(cfg.Postgres.Password).
		Build(ctx)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv := infra.NewServer().
		Host(cfg.Server.Host).
		Port(cfg.Server.Port).
		RouteGroup("/", router.NewRouter(db.Pool()).Register).
		Build()

	go func() {
		slog.Info("server started")
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}

	slog.Info("server stopped")
}
