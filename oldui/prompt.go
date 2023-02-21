package ui

import (
	"strings"

	"github.com/gen64/go-tui"
)

type PromptWidget struct {
	pane *tui.TUIPane
	msg  string
}

func NewPromptWidget(pane *tui.TUIPane) *PromptWidget {
	p := &PromptWidget{
		pane: pane,
		msg:  "",
	}
	pane.SetStyle(tui.NewTUIPaneStyleFrame())
	pane.SetMinHeight(1)
	pane.SetOnDraw(p.update)
	pane.SetOnIterate(p.update)
	return p
}

func (p *PromptWidget) update(pane *tui.TUIPane) int {
	if pane != nil && pane.GetWidth() > 0 {
		printToPane(1, 0, pane, p.msg)
	}
	return 1
}

func (p *PromptWidget) Set(text string) {
	p.msg = text
	if strings.TrimSpace(text) != "" {
		p.pane.Iterate()
	}
}

func (p *PromptWidget) Iterate() {
	p.pane.Iterate()
}
