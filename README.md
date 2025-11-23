# pip-uv

Forces 'pip' commands to run via 'uv pip' transparently.

This tool is designed to be installed in your Python virtual environment. It replaces the standard `pip` executable with a Go binary that transparently redirects all commands to `uv pip`.

## Why?

If you want to ensure that `pip install` (and other pip commands) always use `uv`'s optimized resolution and installation logic—even when running `pip` directly—this shim achieves that with virtually zero overhead.

## Prerequisites

- **Go**: Required to build the shim during installation.
- **uv**: Must be installed and available in your system `PATH`.

## Installation

You can install this package directly with pip. Because it compiles a native Go binary, **you must have Go installed on your system**.

```bash
# Install directly from source
pip install pip-uv
```

Or, if installing from a git repository:

```bash
pip install git+https://github.com/yourusername/pip-uv.git
```

### Manual Installation (No pip)

If you prefer to build and copy the binary manually:

1.  **Build the shim**:
    ```bash
    make build
    ```
2.  **Install into a venv**:
    ```bash
    make install VENV_PATH=/path/to/your/.venv
    ```

## How it works

1.  **Installation**: The `setup.py` script invokes `go build` to compile `main.go` into a binary named `pip`.
2.  **Placement**: Pip installs this binary into your environment's `bin` (or `Scripts`) directory, effectively shadowing the original `pip`.
3.  **Execution**: When you run `pip`, the shim:
    *   Finds the `uv` executable in your `PATH`.
    *   Constructs the command `uv pip [arguments]`.
    *   Executes `uv` using `syscall.Exec`, replacing the shim process entirely (ensuring correct signal handling and exit codes).

## License

 GPL-3.0 license 
