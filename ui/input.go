package ui

import (
	"strings"

	"github.com/gen64/go-tui"
)

type InputWidget struct {
	pane      *tui.TUIPane
	input     string // The unread input buffer
	lastInput string // Last input that was successfully entered
}

func NewInputWidget(pane *tui.TUIPane) *InputWidget {
	i := &InputWidget{
		pane: pane,
	}
	pane.SetStyle(tui.NewTUIPaneStyleFrame())
	pane.SetMinHeight(1)
	return i
}

func (i *InputWidget) onKeyPress(t *tui.TUI, b []byte) (newline bool) {
	// Backspace
	if b[0] == 127 {
		if len(i.input) > 1 {
			i.input = i.input[:len(i.input)-1]
		} else {
			i.input = ""
		}
	} else if string(b) == "\n" {
		i.lastInput = i.input
		i.input = ""
		if strings.HasPrefix(i.lastInput, "print: ") {
			logger.Printf(i.lastInput)
		} else {
			newline = true
		}
	} else {
		i.input += string(b)
	}
	writeLine(1, 0, i.pane, i.input)
	return
}

func (i *InputWidget) Iterate() {
	i.pane.Iterate()
}

// Return the value, clearing it when read
func (i *InputWidget) Value() string {
	ret := i.lastInput
	i.lastInput = ""
	return ret
}

func (i *InputWidget) WaitForInput() string {
	i.lastInput = ""
	for i.lastInput == "" {
		continue
	}
	return i.lastInput
}
