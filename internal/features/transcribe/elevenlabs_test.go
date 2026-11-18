package transcribe

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientTranscribe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/speech-to-text" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("xi-api-key"); got != "secret" {
			t.Errorf("xi-api-key = %q", got)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if got := r.FormValue("model_id"); got != "scribe_v3" {
			t.Errorf("model_id = %q", got)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("form file: %v", err)
		}
		defer file.Close()
		audio, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		if string(audio) != "audio-bytes" {
			t.Errorf("audio = %q", audio)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"language_code":"en","text":"hello there"}`)
	}))
	defer srv.Close()

	c := &client{
		apiKey:  "secret",
		modelID: "scribe_v3",
		url:     srv.URL + "/v1/speech-to-text",
		http:    srv.Client(),
	}

	text, err := c.transcribe(context.Background(), "voice.ogg", strings.NewReader("audio-bytes"))
	if err != nil {
		t.Fatalf("transcribe: %v", err)
	}
	if text != "hello there" {
		t.Fatalf("text = %q, want %q", text, "hello there")
	}
}

func TestClientTranscribeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"detail":{"message":"invalid api key"}}`)
	}))
	defer srv.Close()

	c := &client{
		apiKey:  "secret",
		modelID: "scribe_v3",
		url:     srv.URL,
		http:    srv.Client(),
	}

	_, err := c.transcribe(context.Background(), "voice.ogg", strings.NewReader("audio-bytes"))
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{"401", "invalid api key"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not contain %q", err, want)
		}
	}
}
