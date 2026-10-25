package test

import (
	"code"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatHumanReadable(t *testing.T) {
	require.Equal(t, "0B", code.FormatHumanReadable(0))
	require.Equal(t, "1B", code.FormatHumanReadable(1))
	require.Equal(t, "1023B", code.FormatHumanReadable(1023))
	require.Equal(t, "1.0KB", code.FormatHumanReadable(1024))
	require.Equal(t, "17.6KB", code.FormatHumanReadable(18000))
}

func TestGetPathSize_File(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "a.txt")
	err := os.WriteFile(fpath, []byte("hello world"), 0o644)
	require.NoError(t, err)

	s, err := code.GetPathSize(fpath, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "11B", s)
}

func TestGetPathSize_DirNonRecursive(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "f1"), []byte("aaa"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".hidden"), []byte("bb"), 0o644))

	s, err := code.GetPathSize(dir, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "3B", s)

	sAll, err := code.GetPathSize(dir, false, false, true)
	require.NoError(t, err)
	require.Equal(t, "5B", sAll)
}

func TestGetPathSize_DirRecursive(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "sub"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "sub", "f1"), []byte("abc"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".roothidden"), []byte("d"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "f2"), []byte("ef"), 0o644))

	s, err := code.GetPathSize(dir, true, false, false)
	require.NoError(t, err)
	require.Equal(t, "5B", s)

	sAll, err := code.GetPathSize(dir, true, false, true)
	require.NoError(t, err)
	require.Equal(t, "6B", sAll)
}
