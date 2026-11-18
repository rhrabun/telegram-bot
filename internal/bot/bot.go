package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const onlyVoiceHint = "I only transcribe voice messages right now."

func New(token string, allowed func(int64) bool) (*bot.Bot, error) {
	return bot.New(token,
		bot.WithNotAsyncHandlers(),
		bot.WithMiddlewares(allowList(allowed)),
		bot.WithDefaultHandler(hint),
	)
}

func allowList(allowed func(int64) bool) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message != nil && !allowed(update.Message.Chat.ID) {
				slog.Warn("blocked chat", "chat_id", update.Message.Chat.ID)
				return
			}
			next(ctx, b, update)
		}
	}
}

func hint(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          update.Message.Chat.ID,
		Text:            onlyVoiceHint,
		ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
	}); err != nil {
		slog.Error("send hint", "error", err)
	}
}
