# Publishing to PyPI

This project uses GitHub Actions to automatically build and publish wheels to PyPI when a new tag is pushed.

## Prerequisites

1.  **PyPI Account**: You need an account on [pypi.org](https://pypi.org/).
2.  **Project on PyPI**: You might need to create the project on PyPI first or verify you have permissions.
3.  **Trusted Publishing (Recommended)**:
    *   Go to your project on PyPI (or create it).
    *   Go to **Settings** > **Publishing**.
    *   Add a new **GitHub** publisher.
    *   **Owner**: `guysoft` (replace with your GitHub username).
    *   **Repository**: `pip-uv` (or whatever you named the repo).
    *   **Workflow name**: `release.yml`.
    *   **Environment**: Leave empty or use `pypi` if you configured it.

## How to Release

We have included a helper script `release_script.py` to automate the version bump and tagging process.

1.  **Run the Release Script**:
    ```bash
    ./release_script.py [patch|minor|major]
    ```
    Default is `patch`.

    Example:
    ```bash
    ./release_script.py patch  # 0.1.0 -> 0.1.1
    ./release_script.py minor  # 0.1.0 -> 0.2.0
    ```

    This script will:
    *   Bump the version in `pyproject.toml` (and `setup.py`).
    *   Commit the changes.
    *   Create a git tag (e.g., `v0.1.1`).
    *   Push the commit and tag to GitHub.

2.  **Watch the Build**:
    Go to the "Actions" tab in your GitHub repository. You should see the "Build and Publish" workflow running.
    
    It will:
    *   Build wheels for Windows, macOS (Intel/Apple Silicon), and Linux (x86_64/aarch64).
    *   Build the source distribution.
    *   Upload everything to PyPI.

## Manual Release (Without Script)

1.  **Update Version**:
    Edit `pyproject.toml`.
2.  **Commit and Push**:
    ```bash
    git add .
    git commit -m "Bump version"
    git push
    ```
3.  **Tag**:
    ```bash
    git tag v0.1.1
    git push origin v0.1.1
    ```
