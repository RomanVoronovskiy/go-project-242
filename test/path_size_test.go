package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatHumanReadable(t *testing.T) {
	require.Equal(t, "0B", FormatHumanReadable(0))
	require.Equal(t, "1B", FormatHumanReadable(1))
	require.Equal(t, "1023B", FormatHumanReadable(1023))
	require.Equal(t, "1KB", FormatHumanReadable(1024))
	require.Equal(t, "17.6KB", FormatHumanReadable(18000))
}

func TestGetPathSize_File(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "a.txt")
	err := os.WriteFile(fpath, []byte("hello world"), 0o644)
	require.NoError(t, err)

	s, err := GetPathSize(fpath, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "11B", s)
}

func TestGetPathSize_DirNonRecursive(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "f1"), []byte("aaa"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".hidden"), []byte("bb"), 0o644)

	s, err := GetPathSize(dir, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "3B", s)

	sAll, err := GetPathSize(dir, false, false, true)
	require.NoError(t, err)
	require.Equal(t, "5B", sAll)
}

func TestGetPathSize_DirRecursive(t *testing.T) {
	dir := t.TempDir()
	_ = os.Mkdir(filepath.Join(dir, "sub"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "sub", "f1"), []byte("abc"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".roothidden"), []byte("d"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "f2"), []byte("ef"), 0o644)

	s, err := GetPathSize(dir, true, false, false)
	require.NoError(t, err)
	require.Equal(t, "5B", s)

	sAll, err := GetPathSize(dir, true, false, true)
	require.NoError(t, err)
	require.Equal(t, "6B", sAll)
}
