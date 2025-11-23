# Publishing to PyPI

This project uses GitHub Actions to automatically build and publish wheels to PyPI when a new tag is pushed.

## Prerequisites

1.  **PyPI Account**: You need an account on [pypi.org](https://pypi.org/).
2.  **Project on PyPI**: You might need to create the project on PyPI first or verify you have permissions.
3.  **Trusted Publishing (Recommended)**:
    *   Go to your project on PyPI (or create it).
    *   Go to **Settings** > **Publishing**.
    *   Add a new **GitHub** publisher.
    *   **Owner**: `yourusername` (replace with your GitHub username).
    *   **Repository**: `pip-uv` (or whatever you named the repo).
    *   **Workflow name**: `release.yml`.
    *   **Environment**: Leave empty or use `pypi` if you configured it.

## How to Release

1.  **Update Version**:
    Edit `pyproject.toml` (and `setup.py` if version is hardcoded there) to the new version number (e.g., `0.1.0`).

    ```toml
    [project]
    version = "0.1.0"
    ```

2.  **Commit and Push**:
    ```bash
    git add .
    git commit -m "Bump version to 0.1.0"
    git push
    ```

3.  **Tag the Release**:
    Create a tag starting with `v`.
    ```bash
    git tag v0.1.0
    git push origin v0.1.0
    ```

4.  **Watch the Build**:
    Go to the "Actions" tab in your GitHub repository. You should see the "Build and Publish" workflow running.
    
    It will:
    *   Build wheels for Windows, macOS (Intel/Apple Silicon), and Linux (x86_64/aarch64).
    *   Build the source distribution.
    *   Upload everything to PyPI.

## Manual Publishing (Optional)

If you want to publish manually from your machine:

1.  Install build tools: `pip install build twine`
2.  Build: `python -m build`
3.  Upload: `twine upload dist/*`

