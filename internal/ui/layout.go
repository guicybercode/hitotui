package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Layout struct {
	FileList *FileList
	Preview  *Preview
	Status   string
	Width    int
	Height   int
}

func NewLayout() *Layout {
	return &Layout{
		FileList: NewFileList(),
		Preview:  NewPreview(),
		Status:   "",
		Width:    120,
		Height:   30,
	}
}

func (l *Layout) UpdateSize(width, height int) {
	l.Width = width
	l.Height = height

	fileListWidth := width/2 - 2
	previewWidth := width - fileListWidth - 4

	l.FileList.Width = fileListWidth
	l.FileList.Height = height - 3

	l.Preview.Width = previewWidth
	l.Preview.Height = height - 3
}

func (l *Layout) SetStatus(message string) {
	l.Status = message
}

func (l *Layout) View() string {
	var sections []string

	fileListView := l.FileList.View()
	previewView := l.Preview.View()

	panes := lipgloss.JoinHorizontal(lipgloss.Left, fileListView, previewView)
	sections = append(sections, panes)

	statusBar := l.renderStatusBar()
	sections = append(sections, statusBar)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (l *Layout) renderStatusBar() string {
	selected := l.FileList.GetSelected()
	var info strings.Builder

	if selected != nil {
		if selected.IsDir {
			info.WriteString(DirectoryStyle.Render(selected.Name))
		} else {
			info.WriteString(FileStyle.Render(selected.Name))
			info.WriteString(DimTextStyle.Render(fmt.Sprintf(" (%d bytes)", selected.Size)))
		}
	}

	statusText := l.Status
	if statusText == "" {
		statusText = DimTextStyle.Render("↑↓/j/k: navigate  Enter: open  q: quit  h: toggle hidden")
	}

	left := info.String()
	right := statusText

	width := l.Width - 4
	leftWidth := len(lipgloss.NewStyle().Render(left))
	rightWidth := len(lipgloss.NewStyle().Render(right))

	maxRightLen := width - leftWidth - 2
	if maxRightLen < 0 {
		maxRightLen = 0
	}
	if leftWidth+rightWidth+2 > width && len(right) > maxRightLen {
		if maxRightLen > 3 {
			right = right[:maxRightLen-3] + "..."
		} else {
			right = ""
		}
		rightWidth = len(right)
	}

	paddingLen := width - leftWidth - rightWidth
	if paddingLen < 0 {
		paddingLen = 0
	}
	padding := strings.Repeat(" ", paddingLen)

	statusContent := left + padding + right
	return StatusBarStyle.Width(l.Width).Render(statusContent)
}
