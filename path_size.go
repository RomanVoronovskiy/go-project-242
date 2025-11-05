package code

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("failed to access path: %w", err)
	}
	var totalSize int64
	if fileInfo.Mode().IsRegular() {
		totalSize = fileInfo.Size()
	} else {
		if recursive {
			err = filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !all && strings.HasPrefix(d.Name(), ".") {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				if !d.IsDir() {
					info, err := d.Info()
					if err != nil {
						return err
					}
					totalSize += info.Size()
				}
				return nil
			})
		} else {
			entries, err := os.ReadDir(path)
			if err != nil {
				return "", fmt.Errorf("failed to read directory: %w", err)
			}
			for _, entry := range entries {
				if !all && strings.HasPrefix(entry.Name(), ".") {
					continue
				}
				fullPath := filepath.Join(path, entry.Name())
				fileInfo, err := os.Stat(fullPath)
				if err != nil {
					continue
				}

				if fileInfo.Mode().IsRegular() {
					totalSize += fileInfo.Size()
				}
			}
		}
	}
	if err != nil {
		return "", fmt.Errorf("error during size calculation: %w", err)
	}
	if human {
		return formatHumanReadable(totalSize), nil
	}
	return fmt.Sprintf("%dB", totalSize), nil
}
func formatHumanReadable(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}
