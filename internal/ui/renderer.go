package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Renderer struct {
	out      io.Writer
	theme    Theme
	useColor bool
}

func NewRenderer(out io.Writer, theme Theme, useColor bool) *Renderer {
	return &Renderer{
		out:      out,
		theme:    theme,
		useColor: useColor,
	}
}

func (r *Renderer) Header(text string) {
	r.writeLine(r.style(text, r.theme.Header))
}

func (r *Renderer) Blank() {
	fmt.Fprintln(r.out)
}

func (r *Renderer) Section(title string) {
	r.writeLine(r.style(title, r.theme.SectionTitle))
}

func (r *Renderer) Step(text string) {
	r.writeLine(Indent + StepPrefix + " " + text)
}

func (r *Renderer) StepLog(text string) {
	r.writeLine(Indent + Indent + LogConnector + " " + r.style(text, r.theme.Muted))
}

func (r *Renderer) Result(text string) {
	r.writeLine(Indent + text)
}

func (r *Renderer) TreeLine(prefix, name string) {
	r.writeLine(Indent + prefix + name)
}

func (r *Renderer) style(text string, style lipgloss.Style) string {
	if !r.useColor {
		return text
	}
	return style.Render(text)
}

func (r *Renderer) writeLine(text string) {
	fmt.Fprintln(r.out, strings.TrimRight(text, "\n"))
}
