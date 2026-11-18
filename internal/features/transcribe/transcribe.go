package transcribe

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/rhrabun/telegram-bot/internal/feature"
)

const (
	workerCount    = 4
	requestTimeout = 5 * time.Minute
	replyTimeout   = 30 * time.Second
	maxMessageLen  = 4096
	elevenLabsURL  = "https://api.elevenlabs.io/v1/speech-to-text"
)

type Transcribe struct {
	http   *http.Client
	client *client
	bot    *bot.Bot
	pipe   *pipeline
}

var _ feature.Feature = (*Transcribe)(nil)

func New(apiKey, modelID string) *Transcribe {
	hc := &http.Client{Timeout: requestTimeout}
	return &Transcribe{
		http:   hc,
		client: &client{apiKey: apiKey, modelID: modelID, url: elevenLabsURL, http: hc},
	}
}

func (t *Transcribe) Name() string { return "transcribe" }

func (t *Transcribe) Register(b *bot.Bot) error {
	t.bot = b
	t.pipe = newPipeline(workerCount, requestTimeout, t.work, t.reply)
	t.pipe.start()
	b.RegisterHandlerMatchFunc(isVoice, t.handle)
	return nil
}

func isVoice(update *models.Update) bool {
	return update.Message != nil && update.Message.Voice != nil
}

func (t *Transcribe) handle(_ context.Context, _ *bot.Bot, update *models.Update) {
	msg := update.Message
	t.pipe.enqueue(job{
		chatID:    msg.Chat.ID,
		messageID: msg.ID,
		fileID:    msg.Voice.FileID,
		fileSize:  msg.Voice.FileSize,
	})
}

func (t *Transcribe) work(ctx context.Context, j *job) (string, error) {
	start := time.Now()
	t.chatAction(ctx, j.chatID)

	file, err := t.bot.GetFile(ctx, &bot.GetFileParams{FileID: j.fileID})
	if err != nil {
		return "", fmt.Errorf("download voice message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.bot.FileDownloadLink(file), nil)
	if err != nil {
		return "", fmt.Errorf("download voice message: %w", err)
	}
	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("download voice message: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download voice message: %s", resp.Status)
	}

	text, err := t.client.transcribe(ctx, "voice.ogg", resp.Body)
	if err != nil {
		return "", err
	}

	slog.Info("transcribed", "chat_id", j.chatID, "bytes", j.fileSize, "took", time.Since(start))
	return text, nil
}

func (t *Transcribe) reply(r result) {
	ctx, cancel := context.WithTimeout(context.Background(), replyTimeout)
	defer cancel()

	text := strings.TrimSpace(r.text)
	switch {
	case r.err != nil:
		slog.Error("transcription failed", "chat_id", r.chatID, "error", r.err)
		text = "Transcription failed: " + r.err.Error()
	case text == "":
		text = "No speech detected."
	}

	for i, chunk := range splitText(text, maxMessageLen) {
		params := &bot.SendMessageParams{ChatID: r.chatID, Text: chunk}
		if i == 0 {
			params.ReplyParameters = &models.ReplyParameters{MessageID: r.messageID}
		}
		if _, err := t.bot.SendMessage(ctx, params); err != nil {
			slog.Error("send transcription", "chat_id", r.chatID, "error", err)
			return
		}
	}
}

func (t *Transcribe) chatAction(ctx context.Context, chatID int64) {
	if _, err := t.bot.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: chatID,
		Action: models.ChatActionTyping,
	}); err != nil {
		slog.Warn("send chat action", "error", err)
	}
}
