package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// FileSystem interface to allow mocking in tests
type FileSystem interface {
	Executable() (string, error)
	EvalSymlinks(path string) (string, error)
	ReadFile(name string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
}

// RealFileSystem implements FileSystem using os package
type RealFileSystem struct{}

func (fs RealFileSystem) Executable() (string, error) {
	return os.Executable()
}
func (fs RealFileSystem) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
func (fs RealFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
func (fs RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// CommandResult holds the final command execution plan
type CommandResult struct {
	Bin      string
	Args     []string
	WorkDir  string
	UseShell bool // Not used currently, but good for future
}

// resolveCommand determines the command to run based on args and environment
func resolveCommand(args []string, fs FileSystem) (CommandResult, error) {
	cmdArgs := args

	// Default: uv pip ...
	res := CommandResult{
		Bin:  "uv",
		Args: append([]string{"uv", "pip"}, cmdArgs...),
	}

	// Strategy:
	// 1. Command must be 'install'
	// 2. We must be in a uv-managed project environment (detected via pyvenv.cfg and uv.lock/pyproject.toml)
	// 3. Arguments must be simple packages (no flags starting with -)
	
	if len(cmdArgs) > 0 && cmdArgs[0] == "install" && len(cmdArgs) > 1 {
		exe, err := fs.Executable()
		if err == nil {
			exe, _ = fs.EvalSymlinks(exe)
			
			// Layout: /path/to/project/.venv/bin/pip
			venvBin := filepath.Dir(exe)
			venvRoot := filepath.Dir(venvBin)
			candidateProjectRoot := filepath.Dir(venvRoot)

			// Check 1: Does pyvenv.cfg exist and mention uv?
			hasUvInCfg := false
			if cfgBytes, err := fs.ReadFile(filepath.Join(venvRoot, "pyvenv.cfg")); err == nil {
				if strings.Contains(string(cfgBytes), "uv =") {
					hasUvInCfg = true
				}
			}

			// Check 2: Does project root have uv.lock or pyproject.toml?
			hasLock := false
			if _, err := fs.Stat(filepath.Join(candidateProjectRoot, "uv.lock")); err == nil {
				hasLock = true
			} else if _, err := fs.Stat(filepath.Join(candidateProjectRoot, "pyproject.toml")); err == nil {
				hasLock = true 
			}

			if hasUvInCfg && hasLock {
				// Check 3: Check for simple arguments
				simpleInstall := true
				installArgs := cmdArgs[1:]
				for _, arg := range installArgs {
					if strings.HasPrefix(arg, "-") {
						simpleInstall = false
						break
					}
					if arg == "." {
						simpleInstall = false
					}
				}

				if simpleInstall {
					// Switch to 'uv add'
					res.Args = append([]string{"uv", "add"}, cmdArgs[1:]...)
					res.WorkDir = candidateProjectRoot
				}
			}
		}
	}
	return res, nil
}

func main() {
	// 1. Find where 'uv' is installed on the system
	uvPath, err := exec.LookPath("uv")
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error: 'uv' is not found in your PATH. Please install uv first.")
		os.Exit(1)
	}

	cmd, _ := resolveCommand(os.Args[1:], RealFileSystem{})

	if cmd.WorkDir != "" {
		fmt.Fprintf(os.Stderr, "🚀 Detected uv project at %s\n   Switching 'pip install' to 'uv add'...\n", cmd.WorkDir)
		if err := os.Chdir(cmd.WorkDir); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: Could not change directory to %s: %v. Falling back to pip install.\n", cmd.WorkDir, err)
			// Fallback
			cmd.Args = append([]string{"uv", "pip"}, os.Args[1:]...)
		}
	}

	env := os.Environ()
	if err := syscall.Exec(uvPath, cmd.Args, env); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error executing uv: %v\n", err)
		os.Exit(1)
	}
}
