package frame

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Frame struct {
	winHeight, winWidth int
	margin              *lipgloss.Style
}

func New() *Frame {
	return &Frame{}
}

func (f *Frame) MarginSet(i ...int) *Frame {
	m := lipgloss.NewStyle().Margin(i...)
	f.margin = &m
	return f
}

func (f *Frame) WinSize(msg tea.WindowSizeMsg) {
	f.winHeight = msg.Height
	f.winWidth = msg.Width
}

func (f *Frame) FreeSpace(top, bottom string) int {
	var v int
	if f.margin != nil {
		_, v = f.margin.GetFrameSize()
	}
	return f.winHeight - strings.Count(top, "\n") - strings.Count(bottom, "\n") - 1 - v
}

func (f *Frame) Render(top, bottom string) string {
	freeSpace := f.FreeSpace(top, bottom)

	var padding string
	if freeSpace > 0 {
		padding = strings.Repeat("\n", freeSpace)
	}

	result := top + padding + bottom
	if f.margin != nil {
		return f.margin.Copy().Render(result)
	}
	return result
}

func (f *Frame) Width() int {
	return f.winWidth - f.H()
}

func (f *Frame) Height() int {
	return f.winHeight - f.V()
}

func (f *Frame) WidthFull() int {
	return f.winWidth
}

func (f *Frame) HeightFull() int {
	return f.winHeight
}

func (f *Frame) V() int {
	if f.margin == nil {
		return 0
	}
	_, v := f.margin.GetFrameSize()
	return v
}

func (f *Frame) H() int {
	if f.margin == nil {
		return 0
	}
	_, v := f.margin.GetFrameSize()
	return v
}

func (f *Frame) SingleHeader(header string) string {
	padding := "    "
	wrapL := 38 - 2 - len(padding) - len(header)

	var wrap string
	if wrapL > 0 {
		wrap = strings.Repeat(" ", wrapL)
	}
	return "│" + padding + header + wrap + "│"
}
