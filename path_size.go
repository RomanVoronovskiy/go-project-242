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
		return "", fmt.Errorf("failed to access path: %w", err)
	}

	var total int64
	if info.Mode().IsRegular() {
		total = info.Size()
	} else {
		if recursive {
			total, err = dirSizeRecursive(path, all)
			if err != nil {
				return "", fmt.Errorf("failed to calculate recursive size: %w", err)
			}
		} else {
			total, err = dirSizeNonRecursive(path, all)
			if err != nil {
				return "", fmt.Errorf("failed to calculate directory size: %w", err)
			}
		}
	}

	if human {
		return FormatHumanReadable(total), nil
	}
	return fmt.Sprintf("%dB", total), nil
}

func dirSizeRecursive(root string, all bool) (int64, error) {
	var total int64
	err := filepath.WalkDir(root, func(curr string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
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
			total += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return total, nil
}

func dirSizeNonRecursive(path string, all bool) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, e := range entries {
		if !all && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		full := filepath.Join(path, e.Name())
		info, err := os.Lstat(full)
		if err != nil {
			return 0, fmt.Errorf("failed to stat %q: %w", full, err)
		}
		if info.Mode().IsRegular() {
			total += info.Size()
		}
	}
	return total, nil
}