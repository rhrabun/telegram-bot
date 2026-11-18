package transcribe

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSplitText(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		limit int
		want  []string
	}{
		{name: "short", text: "hello", limit: 10, want: []string{"hello"}},
		{name: "empty", text: "", limit: 10},
		{name: "exact", text: "0123456789", limit: 10, want: []string{"0123456789"}},
		{name: "splits on newline", text: "abc\ndef", limit: 4, want: []string{"abc", "def"}},
		{name: "splits mid word", text: "abcdefgh", limit: 3, want: []string{"abc", "def", "gh"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitText(tt.text, tt.limit)
			if strings.Join(got, "|") != strings.Join(tt.want, "|") {
				t.Fatalf("splitText(%q, %d) = %q, want %q", tt.text, tt.limit, got, tt.want)
			}
		})
	}
}

func TestSplitTextKeepsRunesWhole(t *testing.T) {
	text := strings.Repeat("é", 10)

	chunks := splitText(text, 5)
	for _, chunk := range chunks {
		if !utf8.ValidString(chunk) {
			t.Fatalf("chunk %q is not valid utf-8", chunk)
		}
		if len(chunk) > 5 {
			t.Fatalf("chunk %q is %d bytes, limit is 5", chunk, len(chunk))
		}
	}
	if joined := strings.Join(chunks, ""); joined != text {
		t.Fatalf("joined chunks = %q, want %q", joined, text)
	}
}
