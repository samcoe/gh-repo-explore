package explore

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	viewportStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.NormalBorder())
	}()
)

type fileView struct {
	active   bool
	content  string
	ready    bool
	viewport viewport.Model
}

func NewFileView() *fileView {
	return &fileView{}
}

func (f *fileView) SetSize(width, height int) {
	if !f.ready {
		f.viewport = viewport.New(width, height)
		f.viewport.YPosition = 10
		f.viewport.Style = viewportStyle
		return
	}
	f.viewport.Width = width
	f.viewport.Height = height
}

func (f *fileView) Init() tea.Cmd {
	return nil
}

func (f *fileView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Handle keyboard and mouse events in the viewport
	f.viewport, cmd = f.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return f, tea.Batch(cmds...)
}

func (f *fileView) View() string {
	return f.viewport.View()
}
