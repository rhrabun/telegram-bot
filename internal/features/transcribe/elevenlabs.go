package transcribe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type client struct {
	apiKey  string
	modelID string
	url     string
	http    *http.Client
}

type transcript struct {
	LanguageCode string `json:"language_code"`
	Text         string `json:"text"`
}

func (c *client) transcribe(ctx context.Context, filename string, audio io.Reader) (string, error) {
	var body bytes.Buffer
	form := multipart.NewWriter(&body)

	part, err := form.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	if _, err := io.Copy(part, audio); err != nil {
		return "", fmt.Errorf("read audio: %w", err)
	}
	if err := form.WriteField("model_id", c.modelID); err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	if err := form.Close(); err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, &body)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.Header.Set("xi-api-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("transcribe: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("transcribe: %s: %s", resp.Status, truncate(strings.TrimSpace(string(raw)), 500))
	}

	var out transcript
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return out.Text, nil
}

func truncate(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit]) + "…"
}
