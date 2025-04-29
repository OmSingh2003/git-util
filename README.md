# git-util

[![Release](https://img.shields.io/github/v/release/OmSingh2003/git-util)](https://github.com/OmSingh2003/git-util/releases/latest)
A command-line utility tool written in Go to simplify common Git operations and assist with DevOps workflows.

This tool is developed by [OmSingh2003](https://github.com/OmSingh2003).

## Purpose
efeda
The goal of `git-util` is to provide helpful, automated commands for tasks frequently performed with Git, making repository management easier and faster.

## Features (v0.1.0)

* **Branch Cleaner (`git-util` root command):** Finds and optionally deletes locally merged branches (`-d` to delete, `-n` for dry-run, `-m` to specify main branch).
* **Multi-Repo Status (`status` subcommand):** Checks status (dirty, ahead/behind) of multiple repos in a directory (`-D` to specify directory).
* **Multi-Repo Sync (`sync` subcommand):** Fetches (`-a fetch`) or pulls (`-a pull`) updates across multiple repos (`-D` to specify directory).

## Installation

### Homebrew (Recommended for macOS/Linux)

1.  Tap the repository:
    ```bash
    brew tap OmSingh2003/git-util
    ```
2.  Install `git-util`:
    ```bash
    brew install git-util
    ```
    To upgrade later: `brew upgrade git-util`

### go install

```bash
go install [github.com/OmSingh2003/git-util@latest](https://github.com/OmSingh2003/git-util@latest)
```
Ensure `$HOME/go/bin` is in your `PATH`.

### Manual Download

Download the pre-compiled binary for your operating system from the [GitHub Releases page](https://github.com/OmSingh2003/git-util/releases/latest), extract the archive, and place the `git-util` binary in your desired location (preferably a directory in your `PATH`).

## Usage

### Branch Cleaner (Root Command)

* List potentially deletable merged branches (merged into detected `main`/`master`):
    ```bash
    git-util
    # Or specify main branch:
    git-util -m develop
    ```
* Dry run deletion:
    ```bash
    git-util -d -n
    # Or git-util --delete --dry-run
    ```
* Actually delete merged branches:
    ```bash
    git-util -d
    # Or git-util --delete
    ```

### Multi-Repo Status (`status` subcommand)

* Check status of repos in the current directory:
    ```bash
    git-util status
    ```
* Check status of repos in a specific directory:
    ```bash
    git-util status -D /path/to/your/projects
    # Or git-util status --directory /path/to/your/projects
    ```

### Multi-Repo Sync (`sync` subcommand)

* Fetch updates (`Workspace --prune`) for repos in the current directory (default action):
    ```bash
    git-util sync
    # Or explicitly:
    git-util sync -a fetch
    ```
* Pull updates (`pull --ff-only`) for repos in the current directory:
    ```bash
    git-util sync -a pull
    ```
* Specify directory and action:
    ```bash
    git-util sync -D /path/to/projects -a fetch
    ```

## Development

Clone the repository and build using standard Go commands:

```bash
git clone [https://github.com/OmSingh2003/git-util.git](https://github.com/OmSingh2003/git-util.git)
cd git-util
go build
```

Contributions are welcome! Please open an issue or pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.