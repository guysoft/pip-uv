package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// 1. Find where 'uv' is installed on the system
	// LookPath searches for the executable named file in the directories named by the PATH environment variable.
	uvPath, err := exec.LookPath("uv")
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error: 'uv' is not found in your PATH. Please install uv first.")
		os.Exit(1)
	}

	// 2. Construct the new command: uv pip [args...]
	// os.Args[1:] contains the arguments passed to this program (excluding the program name itself)
	// We want to run: uv pip <original_args>
	// argv[0] for the new process should be "uv" (or the path)
	args := append([]string{"uv", "pip"}, os.Args[1:]...)

	// 3. Execute 'uv', replacing the current process
	// syscall.Exec(argv0, argv, envv)
	// This replaces the current process image with the new one.
	// This is efficient and handles signal propagation/exit codes automatically since it IS the process.
	env := os.Environ()
	if err := syscall.Exec(uvPath, args, env); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error executing uv: %v\n", err)
		os.Exit(1)
	}
}

