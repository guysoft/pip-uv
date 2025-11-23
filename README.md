# pip-uv 🚀

**Tired of forgetting to type `uv pip install`?** 

You know `uv` is faster, better, and stronger, but muscle memory is hard to break. You keep typing `pip install` and waiting... and waiting.

**`pip-uv` is here to save you.**

This package replaces your environment's `pip` command with a lightning-fast shim that automatically redirects everything to `uv pip`. 

Type `pip`, get `uv`. It's that simple.

## Quick Start

Run this **once** in your virtual environment:

```bash
uv pip install pip-uv
```

*(Or just `pip install pip-uv` if you haven't switched yet)*

That's it! Now try it out:

```bash
pip install requests
# 🎉 Actually runs: uv pip install requests
```

## How it works

When you install `pip-uv`, it places a small, optimized binary named `pip` into your virtual environment's `bin` folder. This binary shadows the standard python `pip`.

1.  You type `pip install ...`
2.  The shim intercepts the call.
3.  It instantly replaces itself with `uv pip install ...`
4.  You enjoy pure speed.

## Features

*   **Zero Overhead**: Written in Go, the shim uses `syscall.Exec` to replace the process. No python startup cost.
*   **Transparent**: Passes all arguments and flags directly to `uv`.
*   **Pre-compiled**: Installs instantly on Linux, macOS, and Windows (no Go compiler needed for standard platforms!).

## Prerequisites

*   [**uv**](https://github.com/astral-sh/uv) must be installed and available in your system `PATH`.
    *   *If `uv` isn't found, the shim will complain and exit.*

## License

MIT
