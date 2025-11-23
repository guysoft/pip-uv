package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// MockFileSystem for testing
type MockFileSystem struct {
	ExePath       string
	Files         map[string]string // content of files
	ExistingFiles map[string]bool   // existence of files (for Stat)
}

func (m MockFileSystem) Executable() (string, error) {
	return m.ExePath, nil
}
func (m MockFileSystem) EvalSymlinks(path string) (string, error) {
	return path, nil // Simplification
}
func (m MockFileSystem) ReadFile(name string) ([]byte, error) {
	content, ok := m.Files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return []byte(content), nil
}
func (m MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if m.ExistingFiles[name] {
		return &MockFileInfo{name: filepath.Base(name)}, nil
	}
	// Check if it's in Files map too
	if _, ok := m.Files[name]; ok {
		return &MockFileInfo{name: filepath.Base(name)}, nil
	}
	return nil, os.ErrNotExist
}

// MockFileInfo implements os.FileInfo
type MockFileInfo struct {
	name string
}

func (m *MockFileInfo) Name() string       { return m.name }
func (m *MockFileInfo) Size() int64        { return 0 }
func (m *MockFileInfo) Mode() os.FileMode  { return 0 }
func (m *MockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *MockFileInfo) IsDir() bool        { return false }
func (m *MockFileInfo) Sys() interface{}   { return nil }

func TestResolveCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		exePath  string
		files    map[string]string
		exists   map[string]bool
		expected []string // expected args
		workDir  string   // expected workdir
	}{
		{
			name:     "Standard pip install",
			args:     []string{"install", "requests"},
			exePath:  "/tmp/env/bin/pip",
			files:    map[string]string{},
			expected: []string{"uv", "pip", "install", "requests"},
			workDir:  "",
		},
		{
			name:    "UV Project Detection",
			args:    []string{"install", "requests"},
			exePath: "/app/project/.venv/bin/pip",
			files: map[string]string{
				"/app/project/.venv/pyvenv.cfg": "home = ...\nuv = 0.4.0\n",
			},
			exists: map[string]bool{
				"/app/project/uv.lock": true,
			},
			expected: []string{"uv", "add", "requests"},
			workDir:  "/app/project",
		},
		{
			name:    "UV Project but complex args (fallback)",
			args:    []string{"install", "-r", "requirements.txt"},
			exePath: "/app/project/.venv/bin/pip",
			files: map[string]string{
				"/app/project/.venv/pyvenv.cfg": "uv = 0.4.0",
			},
			exists: map[string]bool{
				"/app/project/uv.lock": true,
			},
			expected: []string{"uv", "pip", "install", "-r", "requirements.txt"},
			workDir:  "",
		},
		{
			name:    "UV Project but install . (fallback)",
			args:    []string{"install", "."},
			exePath: "/app/project/.venv/bin/pip",
			files: map[string]string{
				"/app/project/.venv/pyvenv.cfg": "uv = 0.4.0",
			},
			exists: map[string]bool{
				"/app/project/pyproject.toml": true,
			},
			expected: []string{"uv", "pip", "install", "."},
			workDir:  "",
		},
		{
			name:    "Not UV venv (missing config)",
			args:    []string{"install", "requests"},
			exePath: "/app/project/.venv/bin/pip",
			files:   map[string]string{}, // Empty cfg
			exists: map[string]bool{
				"/app/project/uv.lock": true,
			},
			expected: []string{"uv", "pip", "install", "requests"},
			workDir:  "",
		},
		{
			name:    "UV venv but not project (missing lock)",
			args:    []string{"install", "requests"},
			exePath: "/tmp/random/.venv/bin/pip",
			files: map[string]string{
				"/tmp/random/.venv/pyvenv.cfg": "uv = 0.4.0",
			},
			exists:   map[string]bool{},
			expected: []string{"uv", "pip", "install", "requests"},
			workDir:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := MockFileSystem{
				ExePath:       tt.exePath,
				Files:         tt.files,
				ExistingFiles: tt.exists,
			}
			res, _ := resolveCommand(tt.args, fs)

			if !reflect.DeepEqual(res.Args, tt.expected) {
				t.Errorf("Args = %v, want %v", res.Args, tt.expected)
			}
			if res.WorkDir != tt.workDir {
				t.Errorf("WorkDir = %v, want %v", res.WorkDir, tt.workDir)
			}
		})
	}
}

