package fs

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

type FileInfo struct {
	Name    string
	Path    string
	Size    int64
	IsDir   bool
	ModTime time.Time
}

func ReadDirectory(path string, showHidden bool) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		if !showHidden && entry.Name()[0] == '.' {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		fileInfo := FileInfo{
			Name:    entry.Name(),
			Path:    fullPath,
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime(),
		}
		files = append(files, fileInfo)
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir && !files[j].IsDir {
			return true
		}
		if !files[i].IsDir && files[j].IsDir {
			return false
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

func GetFileInfo(path string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Name:    filepath.Base(path),
		Path:    path,
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
	}, nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ReadFilePreview(path string, maxBytes int64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, maxBytes)
	n, err := file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	return buffer[:n], nil
}

func GetParentDir(path string) string {
	return filepath.Dir(path)
}

func AbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
