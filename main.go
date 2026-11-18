package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	telegram "github.com/rhrabun/telegram-bot/internal/bot"
	"github.com/rhrabun/telegram-bot/internal/config"
	"github.com/rhrabun/telegram-bot/internal/feature"
	"github.com/rhrabun/telegram-bot/internal/features/transcribe"
)

func main() {
	if err := run(); err != nil {
		slog.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	b, err := telegram.New(cfg.TelegramToken, cfg.ChatAllowed)
	if err != nil {
		return fmt.Errorf("telegram bot: %w", err)
	}

	features := []feature.Feature{
		transcribe.New(cfg.ElevenLabsAPIKey, cfg.ElevenLabsModel),
	}
	for _, f := range features {
		if err := f.Register(b); err != nil {
			return fmt.Errorf("register %s: %w", f.Name(), err)
		}
	}

	slog.Info("bot started", "features", len(features), "model", cfg.ElevenLabsModel)
	b.Start(ctx)
	return nil
}
