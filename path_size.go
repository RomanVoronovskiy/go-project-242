package code

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("access %q: %w", path, err)
	}

	size := info.Size()
	if info.IsDir() {
		size, err = directorySize(path, recursive, all)
		if err != nil {
			return "", err
		}
	}

	if human {
		return FormatHumanReadable(size), nil
	}
	return fmt.Sprintf("%dB", size), nil
}

func directorySize(path string, recursive, all bool) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("read directory %q: %w", path, err)
	}

	var total int64
	for _, entry := range entries {
		if !all && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		size, err := entrySize(path, entry, recursive, all)
		if err != nil {
			return 0, err
		}
		total += size
	}
	return total, nil
}

func entrySize(parent string, entry fs.DirEntry, recursive, all bool) (int64, error) {
	path := filepath.Join(parent, entry.Name())
	if entry.IsDir() {
		if recursive {
			return directorySize(path, recursive, all)
		}
		return 0, nil
	}

	info, err := entry.Info()
	if err != nil {
		return 0, fmt.Errorf("access %q: %w", path, err)
	}
	return info.Size(), nil
}
