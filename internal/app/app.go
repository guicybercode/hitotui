package app

import (
	"fmt"
	"os"

	"hitotui/internal/fs"
	"hitotui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type model struct {
	layout      *ui.Layout
	currentPath string
	showHidden  bool
	status      string
	width       int
	height      int
	quitting    bool
}

type statusMsg string
type errMsg error

func initialModel() model {
	path, _ := os.Getwd()
	if len(os.Args) > 1 {
		if fs.Exists(os.Args[1]) {
			path = fs.AbsPath(os.Args[1])
		}
	}

	m := model{
		layout:      ui.NewLayout(),
		currentPath: path,
		showHidden:  false,
		status:      "",
		width:       120,
		height:      30,
	}

	m.loadDirectory()
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout.UpdateSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			m.layout.FileList.MoveUp()
			m.updatePreview()
			return m, nil

		case "down", "j":
			m.layout.FileList.MoveDown()
			m.updatePreview()
			return m, nil

		case "enter":
			selected := m.layout.FileList.GetSelected()
			if selected != nil {
				if selected.Name == ".." {
					parent := fs.GetParentDir(m.currentPath)
					currentAbs := fs.AbsPath(m.currentPath)
					parentAbs := fs.AbsPath(parent)
					if currentAbs != parentAbs {
						m.currentPath = parent
						m.loadDirectory()
					}
				} else if selected.IsDir {
					m.currentPath = selected.Path
					m.loadDirectory()
				}
			}
			return m, nil

		case "backspace", "h":
			m.goToParent()
			return m, nil

		case ".":
			m.showHidden = !m.showHidden
			m.loadDirectory()
			m.setStatus(fmt.Sprintf("Hidden files: %v", m.showHidden))
			return m, nil

		case "c":
			return m, m.copyFile()

		case "x":
			return m, m.cutFile()

		case "d":
			return m, m.deleteFile()

		case "r":
			return m, m.renameFile()

		case "n":
			return m, m.createDirectory()

		case "esc":
			m.status = ""
			return m, nil

		case "/":
			return m, nil
		}

	case statusMsg:
		m.status = string(msg)
		return m, nil

	case errMsg:
		m.setStatus(fmt.Sprintf("Error: %v", msg))
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return m.layout.View()
}

func (m *model) loadDirectory() {
	files, err := fs.ReadDirectory(m.currentPath, m.showHidden)
	if err != nil {
		m.setStatus(fmt.Sprintf("Error reading directory: %v", err))
		return
	}

	parent := fs.GetParentDir(m.currentPath)
	currentAbs := fs.AbsPath(m.currentPath)
	parentAbs := fs.AbsPath(parent)

	if currentAbs != parentAbs {
		parentFile := fs.FileInfo{
			Name:  "..",
			Path:  parent,
			IsDir: true,
		}
		files = append([]fs.FileInfo{parentFile}, files...)
	}

	m.layout.FileList.SetFiles(files)
	m.layout.FileList.ShowHidden = m.showHidden
	m.updatePreview()
	m.currentPath = currentAbs
}

func (m *model) updatePreview() {
	selected := m.layout.FileList.GetSelected()
	m.layout.Preview.SetFile(selected)
}

func (m *model) setStatus(msg string) {
	m.status = msg
	m.layout.SetStatus(msg)
}

func (m *model) goToParent() {
	parent := fs.GetParentDir(m.currentPath)
	currentAbs := fs.AbsPath(m.currentPath)
	parentAbs := fs.AbsPath(parent)
	if currentAbs != parentAbs {
		m.currentPath = parent
		m.loadDirectory()
		m.layout.FileList.Selected = 0
	}
}

func (m *model) copyFile() tea.Cmd {
	return func() tea.Msg {
		selected := m.layout.FileList.GetSelected()
		if selected == nil || selected.IsDir {
			return statusMsg("Cannot copy directory")
		}
		return statusMsg(fmt.Sprintf("Copy mode: %s (press Enter on destination)", selected.Name))
	}
}

func (m *model) cutFile() tea.Cmd {
	return func() tea.Msg {
		selected := m.layout.FileList.GetSelected()
		if selected == nil || selected.IsDir {
			return statusMsg("Cannot move directory")
		}
		return statusMsg(fmt.Sprintf("Cut mode: %s (press Enter on destination)", selected.Name))
	}
}

func (m *model) deleteFile() tea.Cmd {
	return func() tea.Msg {
		selected := m.layout.FileList.GetSelected()
		if selected == nil {
			return errMsg(fmt.Errorf("no file selected"))
		}
		return statusMsg(fmt.Sprintf("Delete: %s (not implemented - use system commands for safety)", selected.Name))
	}
}

func (m *model) renameFile() tea.Cmd {
	return func() tea.Msg {
		selected := m.layout.FileList.GetSelected()
		if selected == nil {
			return errMsg(fmt.Errorf("no file selected"))
		}
		return statusMsg(fmt.Sprintf("Rename: %s (not implemented)", selected.Name))
	}
}

func (m *model) createDirectory() tea.Cmd {
	return func() tea.Msg {
		return statusMsg("Create directory (not implemented)")
	}
}

func Run() error {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return fmt.Errorf("not a terminal")
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
