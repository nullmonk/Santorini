package ui

import (
	"strings"

	"github.com/gen64/go-tui"
)

func printToPane(x, y int, p *tui.TUIPane, text string) {
	for i, line := range wordWrap(text, p.GetWidth(), "↵") {
		writeLine(x, y+i, p, line)
	}
}

// Write a line to the given pane, clearing the underlying content if needed
func writeLine(x, y int, p *tui.TUIPane, line string) {
	style := p.GetStyle()
	borderWidth := 0
	if style != nil && len(strings.TrimSpace(style.E+style.W)) > 0 {
		borderWidth = 1
	}
	maxWidth := p.GetWidth() - x - borderWidth - borderWidth - x
	if len(line) > maxWidth {
		p.Write(x, y, line[len(line)-maxWidth:], false)
	} else {
		p.Write(x, y, line+strings.Repeat(" ", maxWidth-len(line)), false)
	}
}

func wordWrap(text string, lineWidth int, suffix ...string) []string {
	rows := make([]string, 0, 2)
	suf := ""
	if len(suffix) > 0 {
		suf = suffix[0]
		lineWidth -= len(suf)
	}
	size := len(text)
	for size > lineWidth {
		rows = append(rows, strings.TrimSpace(text[:lineWidth])+suf)
		text = strings.TrimSpace(text[lineWidth:])
		size = len(text)
	}
	rows = append(rows, text)
	return rows
}
