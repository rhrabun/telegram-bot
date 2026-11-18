package transcribe

import (
	"strings"
	"unicode/utf8"
)

func splitText(text string, limit int) []string {
	var chunks []string

	for len(text) > limit {
		cut := strings.LastIndexByte(text[:limit], '\n')
		if cut <= 0 {
			cut = limit
			for cut > 0 && !utf8.RuneStart(text[cut]) {
				cut--
			}
			if cut == 0 {
				cut = limit
			}
		}
		chunks = append(chunks, strings.TrimRight(text[:cut], "\n"))
		text = strings.TrimLeft(text[cut:], "\n")
	}

	if text != "" {
		chunks = append(chunks, text)
	}
	return chunks
}
