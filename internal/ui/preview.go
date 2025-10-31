package ui

import (
	"fmt"
	"strings"

	"hitotui/internal/fs"
)

type Preview struct {
	File     *fs.FileInfo
	Content  string
	Width    int
	Height   int
	MaxBytes int64
}

func NewPreview() *Preview {
	return &Preview{
		Width:    60,
		Height:   20,
		MaxBytes: 5000,
	}
}

func (p *Preview) SetFile(file *fs.FileInfo) {
	p.File = file
	p.Content = ""
	if file != nil {
		p.loadContent()
	}
}

func (p *Preview) loadContent() {
	if p.File == nil {
		return
	}

	if p.File.IsDir {
		p.Content = p.renderDirectoryInfo()
		return
	}

	content, err := fs.ReadFilePreview(p.File.Path, p.MaxBytes)
	if err != nil {
		p.Content = ErrorTextStyle.Render(fmt.Sprintf("Error: %v", err))
		return
	}

	contentStr := string(content)
	if !isTextFile(contentStr) {
		p.Content = DimTextStyle.Render(fmt.Sprintf("Binary file (%d bytes)\nCannot preview", p.File.Size))
		return
	}

	p.Content = formatContent(contentStr)
}

func (p *Preview) renderDirectoryInfo() string {
	var info strings.Builder
	info.WriteString(DirectoryStyle.Render(fmt.Sprintf("📁 %s\n\n", p.File.Name)))
	info.WriteString(DimTextStyle.Render(fmt.Sprintf("Path: %s\n", p.File.Path)))
	info.WriteString(DimTextStyle.Render(fmt.Sprintf("Modified: %s\n", p.File.ModTime.Format("2006-01-02 15:04:05"))))
	return info.String()
}

func isTextFile(content string) bool {
	for _, r := range content {
		if r < 32 && r != 9 && r != 10 && r != 13 {
			return false
		}
		if r > 127 && r < 160 {
			return false
		}
	}
	return true
}

func formatContent(content string) string {
	lines := strings.Split(content, "\n")
	var formatted []string
	maxLines := 50

	if len(lines) > maxLines {
		formatted = lines[:maxLines]
		formatted = append(formatted, DimTextStyle.Render(fmt.Sprintf("\n... (%d more lines)", len(lines)-maxLines)))
	} else {
		formatted = lines
	}

	limited := make([]string, 0)
	for _, line := range formatted {
		if len(line) > 100 {
			limited = append(limited, line[:100]+DimTextStyle.Render("..."))
		} else {
			limited = append(limited, line)
		}
	}

	return strings.Join(limited, "\n")
}

func (p *Preview) View() string {
	if p.File == nil {
		return DimTextStyle.Render("No file selected")
	}

	header := fmt.Sprintf("%s", p.File.Name)
	if p.File.IsDir {
		header = DirectoryStyle.Render(header)
	} else {
		header = FileStyle.Render(header)
		header += DimTextStyle.Render(fmt.Sprintf(" (%d bytes)", p.File.Size))
	}

	content := p.Content
	if p.Content == "" {
		content = DimTextStyle.Render("Loading...")
	}

	view := fmt.Sprintf("%s\n\n%s", header, content)
	return PreviewStyle.Width(p.Width).Height(p.Height).Render(view)
}
