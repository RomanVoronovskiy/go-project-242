package code

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize(t *testing.T) {
	root := t.TempDir()
	files := map[string]int{
		"file.txt":              3,
		".hidden":               2,
		"sub/file.txt":          5,
		"sub/.hidden":           7,
		"sub/deep/file.txt":     13,
		".secret/file.txt":      11,
		".secret/deep/file.txt": 17,
	}
	for name, size := range files {
		path := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, make([]byte, size), 0o644))
	}

	tests := []struct {
		name           string
		recursive, all bool
		want           string
	}{
		{"first level", false, false, "3B"},
		{"first level including hidden", false, true, "5B"},
		{"recursive", true, false, "21B"},
		{"recursive including hidden", true, true, "58B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, human := range []bool{false, true} {
				size, err := GetPathSize(root, tt.recursive, human, tt.all)
				require.NoError(t, err)
				require.Equal(t, tt.want, size)
			}
		})
	}
}

func TestGetPathSizeFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "output.dat")
	file, err := os.Create(path)
	require.NoError(t, err)
	require.NoError(t, file.Truncate(25165824))
	require.NoError(t, file.Close())
	for _, recursive := range []bool{false, true} {
		for _, all := range []bool{false, true} {
			for _, human := range []bool{false, true} {
				want := "25165824B"
				if human {
					want = "24.0MB"
				}
				size, err := GetPathSize(path, recursive, human, all)
				require.NoError(t, err)
				require.Equal(t, want, size)
			}
		}
	}
}

func TestGetPathSizeEmpty(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "empty")
	require.NoError(t, os.WriteFile(path, nil, 0o644))
	for _, path := range []string{root, path} {
		size, err := GetPathSize(path, true, true, true)
		require.NoError(t, err)
		require.Equal(t, "0B", size)
	}
}

func TestGetPathSizeExplicitHiddenRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".root")
	require.NoError(t, os.Mkdir(root, 0o755))
	path := filepath.Join(root, ".hidden")
	require.NoError(t, os.WriteFile(path, []byte("abc"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "visible"), []byte("hello"), 0o644))
	for _, recursive := range []bool{false, true} {
		size, err := GetPathSize(root, recursive, false, false)
		require.NoError(t, err)
		require.Equal(t, "5B", size)
	}
	size, err := GetPathSize(path, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "3B", size)
}

func TestGetPathSizeMissingPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	size, err := GetPathSize(path, true, false, true)
	require.ErrorIs(t, err, fs.ErrNotExist)
	require.Empty(t, size)
	require.Contains(t, err.Error(), path)
}

func TestDirectoryReadError(t *testing.T) {
	_, err := directorySize(filepath.Join(t.TempDir(), "missing"), true, true)
	require.ErrorIs(t, err, fs.ErrNotExist)
}

func TestEntryDisappears(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows directory entries cache file metadata")
	}
	root := t.TempDir()
	path := filepath.Join(root, "file")
	require.NoError(t, os.WriteFile(path, []byte("abc"), 0o644))
	entries, err := os.ReadDir(root)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.NoError(t, os.Remove(path))
	_, err = entrySize(root, entries[0], false, true)
	require.ErrorIs(t, err, fs.ErrNotExist)
	require.Contains(t, err.Error(), path)
}

func TestGetPathSizeUnreadableDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits are not enforced on Windows")
	}
	root := t.TempDir()
	sub := filepath.Join(root, "blocked")
	require.NoError(t, os.Mkdir(sub, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(sub, "file"), []byte("abc"), 0o644))
	require.NoError(t, os.Chmod(sub, 0))
	t.Cleanup(func() { require.NoError(t, os.Chmod(sub, 0o755)) })
	if _, err := os.ReadDir(sub); err == nil {
		t.Skip("current user bypasses Unix permission bits")
	}
	for _, path := range []string{root, sub} {
		size, err := GetPathSize(path, true, false, true)
		require.True(t, errors.Is(err, fs.ErrPermission), "expected permission error, got %v", err)
		require.Empty(t, size)
	}
}

func TestGetPathSizeSymlinks(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "file"), []byte("abc"), 0o644))
	links := map[string]string{"link": "file", "broken": "missing", "loop": "."}
	var total int64 = 3
	for name, target := range links {
		path := filepath.Join(root, name)
		if err := os.Symlink(target, path); err != nil {
			if runtime.GOOS == "windows" {
				t.Skipf("symlink creation requires developer mode or elevated rights: %v", err)
			}
			require.NoError(t, err)
		}
		info, err := os.Lstat(path)
		require.NoError(t, err)
		total += info.Size()
		size, err := GetPathSize(path, true, false, true)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%dB", info.Size()), size)
	}
	for _, recursive := range []bool{false, true} {
		size, err := GetPathSize(root, recursive, false, true)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%dB", total), size)
	}
}

func TestFormatHumanReadableUnits(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{0, "0B"}, {1, "1B"}, {123, "123B"}, {1023, "1023B"},
		{1024, "1.0KB"}, {1536, "1.5KB"}, {18000, "17.6KB"},
		{1234567, "1.2MB"}, {1 << 20, "1.0MB"}, {25165824, "24.0MB"},
		{1 << 30, "1.0GB"}, {1 << 40, "1.0TB"}, {1 << 50, "1.0PB"},
		{1 << 60, "1.0EB"}, {math.MaxInt64, "8.0EB"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			require.Equal(t, tt.want, FormatHumanReadable(tt.size))
		})
	}
}
