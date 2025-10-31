package actions

import (
	"fmt"
	"hitotui/internal/fs"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

func DeleteFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func RenameFile(oldPath, newName string) error {
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)
	return os.Rename(oldPath, newPath)
}

func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

func ValidatePath(path string) bool {
	return fs.Exists(path)
}

func FormatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("Error: %v", err)
}
