package ui

import (
	"strings"

	"github.com/gen64/go-tui"
)

type SelectWidget struct {
	pane *tui.TUIPane
	msg  string
}

func NewSelectWidget(pane *tui.TUIPane) *SelectWidget {
	p := &SelectWidget{
		pane: pane,
		msg:  "",
	}
	pane.SetStyle(tui.NewTUIPaneStyleFrame())
	pane.SetMinHeight(10)
	pane.SetOnDraw(p.update)
	pane.SetOnIterate(p.update)
	return p
}

func (p *SelectWidget) update(pane *tui.TUIPane) int {
	if pane != nil && pane.GetWidth() > 0 {
		printToPane(1, 0, pane, p.msg)
	}
	return 1
}

func (p *SelectWidget) Set(text string) {
	p.msg = text
	if strings.TrimSpace(text) != "" {
		p.pane.Iterate()
	}
}

func (p *SelectWidget) Iterate() {
	p.pane.Iterate()
}
