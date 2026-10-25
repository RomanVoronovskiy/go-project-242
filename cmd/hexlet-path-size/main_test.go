package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandOutput(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "file"), make([]byte, 1024), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".hidden"), make([]byte, 512), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(root, "sub"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "sub", "file"), make([]byte, 2048), 0o644))
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default", []string{root}, "1024B"},
		{"long flags before path", []string{"--human", "--all", "--recursive", root}, "3.5KB"},
		{"short flags after path", []string{root, "-H", "-a", "-r"}, "3.5KB"},
		{"mixed positions", []string{"-r", root, "-H"}, "3.0KB"},
		{"explicit false", []string{root, "--human=false", "--all=false", "--recursive=false"}, "1024B"},
		{"recursive only", []string{root, "-r"}, "3072B"},
		{"all only", []string{root, "-a"}, "1536B"},
		{"file", []string{filepath.Join(root, "file"), "-H"}, "1.0KB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			cmd := newCommand()
			cmd.Writer = &output
			cmd.ErrWriter = io.Discard
			err := cmd.Run(context.Background(), append([]string{"hexlet-path-size"}, tt.args...))
			require.NoError(t, err)
			path := root
			if tt.name == "file" {
				path = filepath.Join(root, "file")
			}
			require.Equal(t, tt.want+"\t"+path+"\n", output.String())
		})
	}
}

func TestCommandHelp(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		t.Run(flag, func(t *testing.T) {
			var output bytes.Buffer
			cmd := newCommand()
			cmd.Writer = &output
			require.NoError(t, cmd.Run(context.Background(), []string{"hexlet-path-size", flag}))
			for _, text := range []string{
				"hexlet-path-size [global options] <path>",
				"print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
				"--recursive, -r", "recursive size of directories (default: false)",
				"--human, -H", "human-readable sizes (auto-select unit) (default: false)",
				"--all, -a", "include hidden files and directories (default: false)",
				"--help, -h", "show help",
			} {
				require.Contains(t, output.String(), text)
			}
		})
	}
}

func TestCommandErrors(t *testing.T) {
	tests := [][]string{
		{},
		{"first", "second"},
		{filepath.Join(t.TempDir(), "missing")},
		{"--unknown"},
		{"--human=invalid"},
	}
	for _, args := range tests {
		cmd := newCommand()
		cmd.Writer = io.Discard
		cmd.ErrWriter = io.Discard
		require.Error(t, cmd.Run(context.Background(), append([]string{"hexlet-path-size"}, args...)))
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestCommandOutputError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")
	require.NoError(t, os.WriteFile(path, []byte("hello"), 0o644))
	cmd := newCommand()
	cmd.Writer = failingWriter{}
	require.ErrorIs(t, cmd.Run(context.Background(), []string{"hexlet-path-size", path}), io.ErrClosedPipe)
}

func TestCommandPathWithSpaces(t *testing.T) {
	path := filepath.Join(t.TempDir(), "файл с пробелами.txt")
	require.NoError(t, os.WriteFile(path, []byte("hello\n"), 0o644))
	var output bytes.Buffer
	cmd := newCommand()
	cmd.Writer = &output
	require.NoError(t, cmd.Run(context.Background(), []string{"hexlet-path-size", path}))
	require.Equal(t, "6B\t"+path+"\n", output.String())
}

func TestRunExitStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")
	require.NoError(t, os.WriteFile(path, []byte("hello\n"), 0o644))
	var stdout, stderr bytes.Buffer
	require.Equal(t, 0, run([]string{"hexlet-path-size", path}, &stdout, &stderr))
	require.Equal(t, "6B\t"+path+"\n", stdout.String())
	require.Empty(t, stderr.String())

	stdout.Reset()
	require.Equal(t, 1, run([]string{"hexlet-path-size", path + ".missing"}, &stdout, &stderr))
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), path+".missing")
}
