// Package frame used for frame utils TUI application.
package frame

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Frame struct with frame configuration.
type Frame struct {
	margin              *lipgloss.Style
	winHeight, winWidth int
}

// New creates new frame.
func New() *Frame {
	return &Frame{}
}

// MarginSet set margins for TUI application.
func (f *Frame) MarginSet(i ...int) *Frame {
	m := lipgloss.NewStyle().Margin(i...)
	f.margin = &m
	return f
}

// WinSize sets windows size using object tea.WindowSizeMsg. When window size changes usualy creates object tea.WindowSizeMsg.
func (f *Frame) WinSize(msg tea.WindowSizeMsg) {
	f.winHeight = msg.Height
	f.winWidth = msg.Width
}

// FreeSpace - according to passed arguments and window height returns free space.
func (f *Frame) FreeSpace(top, bottom string) int {
	var v int
	if f.margin != nil {
		_, v = f.margin.GetFrameSize()
	}
	return f.winHeight - strings.Count(top, "\n") - strings.Count(bottom, "\n") - 1 - v
}

// Render - according to window size and margins renders UI.
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

// Width of window without margins.
func (f *Frame) Width() int {
	return f.winWidth - f.H()
}

// Height of window without margins.
func (f *Frame) Height() int {
	return f.winHeight - f.V()
}

// WidthFull - width of window.
func (f *Frame) WidthFull() int {
	return f.winWidth
}

// HeightFull - height of window.
func (f *Frame) HeightFull() int {
	return f.winHeight
}

// V - vertial margin.
func (f *Frame) V() int {
	if f.margin == nil {
		return 0
	}
	_, v := f.margin.GetFrameSize()
	return v
}

// H - horisontal margin.
func (f *Frame) H() int {
	if f.margin == nil {
		return 0
	}
	_, v := f.margin.GetFrameSize()
	return v
}

// SingleHeader - draws text inside frame.
func (f *Frame) SingleHeader(header string) string {
	padding := "    "
	wrapL := 38 - 2 - len(padding) - len(header)

	var wrap string
	if wrapL > 0 {
		wrap = strings.Repeat(" ", wrapL)
	}
	return "│" + padding + header + wrap + "│"
}
