package svg

import (
	"strings"
	"testing"

	assert "github.com/alecthomas/assert/v2"
)

func TestWriteFontStyle(t *testing.T) {
	tests := []struct {
		name     string
		format   FontFormat
		expected string
	}{
		{"WOFF", WOFF, "src: url(data:font/woff;charset=utf-8;base64,AAAA) format('woff');\n"},
		{"WOFF2", WOFF2, "src: url(data:font/woff2;charset=utf-8;base64,AAAA) format('woff2');\n"},
		{"TrueType", TRUETYPE, "src: url(data:font/ttf;charset=utf-8;base64,AAAA) format('truetype');\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := New(EmbedFont("Test", "AAAA", test.format))
			w := &strings.Builder{}
			assert.NoError(t, f.writeFontStyle(w))
			assert.Contains(t, w.String(), test.expected)
		})
	}
}
