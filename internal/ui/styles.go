package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor   = lipgloss.Color("39")
	SecondaryColor = lipgloss.Color("135")
	SuccessColor   = lipgloss.Color("46")
	WarningColor   = lipgloss.Color("220")
	ErrorColor     = lipgloss.Color("196")
	DimColor       = lipgloss.Color("240")
	BorderColor    = lipgloss.Color("99")

	AppStyle = lipgloss.NewStyle().
			Padding(0, 1)

	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	FileListStyle = PaneStyle.Copy().
			BorderForeground(PrimaryColor)

	PreviewStyle = PaneStyle.Copy().
			BorderForeground(SecondaryColor)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(PrimaryColor).
			Bold(true)

	DirectoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	HiddenStyle = lipgloss.NewStyle().
			Foreground(DimColor).
			Italic(true)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	SuccessTextStyle = lipgloss.NewStyle().
				Foreground(SuccessColor)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(DimColor)
)

func GetFileStyle(name string, isDir bool, isHidden bool, isSelected bool) lipgloss.Style {
	if isSelected {
		return SelectedStyle
	}
	if isHidden {
		return HiddenStyle
	}
	if isDir {
		return DirectoryStyle
	}
	return FileStyle
}
