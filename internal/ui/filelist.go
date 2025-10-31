package ui

import (
	"strings"

	"hitotui/internal/fs"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type FileList struct {
	Files      []fs.FileInfo
	Selected   int
	ShowHidden bool
	Width      int
	Height     int
}

func NewFileList() *FileList {
	return &FileList{
		Files:    []fs.FileInfo{},
		Selected: 0,
		Width:    50,
		Height:   20,
	}
}

func (f *FileList) SetFiles(files []fs.FileInfo) {
	f.Files = files
	if f.Selected >= len(f.Files) {
		f.Selected = max(0, len(f.Files)-1)
	}
}

func (f *FileList) MoveUp() {
	if f.Selected > 0 {
		f.Selected--
	}
}

func (f *FileList) MoveDown() {
	if f.Selected < len(f.Files)-1 {
		f.Selected++
	}
}

func (f *FileList) GetSelected() *fs.FileInfo {
	if len(f.Files) == 0 || f.Selected < 0 || f.Selected >= len(f.Files) {
		return nil
	}
	return &f.Files[f.Selected]
}

func (f *FileList) View() string {
	if len(f.Files) == 0 {
		return DimTextStyle.Render("No files")
	}

	var lines []string
	startIdx := 0
	endIdx := len(f.Files)

	if len(f.Files) > f.Height {
		if f.Selected >= f.Height {
			startIdx = f.Selected - f.Height + 1
			endIdx = startIdx + f.Height
			if endIdx > len(f.Files) {
				endIdx = len(f.Files)
				startIdx = endIdx - f.Height
			}
		} else {
			endIdx = f.Height
		}
	}

	for i := startIdx; i < endIdx && i < len(f.Files); i++ {
		file := f.Files[i]
		isSelected := i == f.Selected
		isHidden := len(file.Name) > 0 && file.Name[0] == '.' && file.Name != ".."

		prefix := "  "
		if file.Name == ".." {
			prefix = "⬆️  "
		} else if file.IsDir {
			prefix = "📁 "
		} else {
			prefix = "📄 "
		}

		style := GetFileStyle(file.Name, file.IsDir, isHidden, isSelected)
		name := style.Render(prefix + file.Name)

		if file.IsDir {
			name += DimTextStyle.Render("/")
		}

		lines = append(lines, name)
	}

	content := strings.Join(lines, "\n")
	return FileListStyle.Width(f.Width).Height(f.Height).Render(content)
}
