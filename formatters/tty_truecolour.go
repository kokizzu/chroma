package formatters

import (
	"fmt"
	"io"
	"iter"
	"regexp"

	"github.com/alecthomas/chroma/v3"
)

// TTY16m is a true-colour terminal formatter.
var TTY16m = Register("terminal16m", chroma.FormatterFunc(trueColourFormatter))

var crOrCrLf = regexp.MustCompile(`\r?\n`)

// Print the text with the given formatting, resetting the formatting at the end
// of each line and resuming it on the next line.
//
// This way, a pager (like https://github.com/walles/moar for example) can show
// any line in the output by itself, and it will get the right formatting.
func writeToken(w io.Writer, formatting string, text string) error {
	if formatting == "" {
		_, err := fmt.Fprint(w, text)
		return err
	}

	newlineIndices := crOrCrLf.FindAllStringIndex(text, -1)

	afterLastNewline := 0
	for _, indices := range newlineIndices {
		newlineStart, afterNewline := indices[0], indices[1]
		_, err := fmt.Fprint(w, formatting, text[afterLastNewline:newlineStart], "\033[0m", text[newlineStart:afterNewline])
		if err != nil {
			return err
		}
		afterLastNewline = afterNewline
	}

	if afterLastNewline < len(text) {
		// Print whatever is left after the last newline
		_, err := fmt.Fprint(w, formatting, text[afterLastNewline:], "\033[0m")
		return err
	}
	return nil
}

func trueColourFormatter(w io.Writer, style *chroma.Style, it iter.Seq[chroma.Token]) error {
	style = clearBackground(style)
	for token := range it {
		entry := style.Get(token.Type)
		if entry.IsZero() {
			if _, err := fmt.Fprint(w, token.Value); err != nil {
				return err
			}
			continue
		}

		formatting := attrEscapeSequence(entry)
		if entry.Colour.IsSet() {
			formatting += fmt.Sprintf("\033[38;2;%d;%d;%dm", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue())
		}
		if entry.Background.IsSet() {
			formatting += fmt.Sprintf("\033[48;2;%d;%d;%dm", entry.Background.Red(), entry.Background.Green(), entry.Background.Blue())
		}

		if err := writeToken(w, formatting, token.Value); err != nil {
			return err
		}
	}
	return nil
}
