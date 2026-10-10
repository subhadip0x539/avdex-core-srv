package main

import (
	"log/slog"
	"os"

	"github.com/subhadip0x539/arc-dex-core-srv/src/config"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/app"
)

func main() {
	cfg, err := config.NewConfigBuilder().Load().Build()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	app.Run(cfg)
}
