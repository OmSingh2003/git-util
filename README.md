# git-util

A command-line utility tool written in Go to simplify common Git operations and assist with DevOps workflows.

This tool is developed by [OmSingh2003](https://github.com/OmSingh2003).

## Purpose

The goal of `git-util` is to provide helpful, automated commands for tasks frequently performed with Git, making repository management easier and faster.

## Features

Currently implemented or planned features:

* **Git Branch Cleaner:**
    * Identifies local Git branches that have already been merged into a specified main branch (e.g., `main`, `master`).
    * Provides an option to delete these merged branches safely (`--delete`).
    * Includes a dry-run mode (`--dry-run`) to preview deletions.
    * Allows specifying the main branch to compare against (`--main`).
* *(More features may be added in the future)*

## Installation

You can install `git-util` using `go install`:

```bash
go install [github.com/OmSingh2003/git-util@latest](https://github.com/OmSingh2003/git-util@latest)
```

Make sure your `$GOPATH/bin` or `$HOME/go/bin` directory is in your system's `PATH` environment variable.

## Building from Source

Alternatively, you can build it from the source code:

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/OmSingh2003/git-util.git](https://github.com/OmSingh2003/git-util.git)
    cd git-util
    ```
2.  **Build the binary:** 
    ```bash
    go build
    ```
    This will create the `git-util` executable in the current directory.

## Usage

### Branch Cleaner

This is the initial command implemented directly on the root `git-util` command.

**1. List Merged Branches:**

To see which local branches have already been merged into your main branch (defaults to checking against `main` or `master`), simply run:

```bash
git-util
```

**2. Specify Main Branch:**

If your main development branch is not `main` or `master`, use the `--main` flag:

```bash
git-util --main develop
```

**3. Delete Merged Branches:**

*First, run without `--delete` to see which branches will be targeted.*

To delete the identified merged branches, use the `--delete` flag. **Use with caution!**

```bash
# This will attempt to delete the merged branches
git-util --delete

# Delete branches merged into 'develop'
git-util --main develop --delete
```

**4. Dry Run Deletion:**

To see which branches *would* be deleted without actually performing the deletion, use `--delete` along with `--dry-run`:

```bash
git-util --delete --dry-run

# Dry run targeting 'develop'
git-util --main develop --delete --dry-run
```

## Development

Feel free to contribute or report issues on the [GitHub repository](https://github.com/OmSingh2003/git-util).

Build the project using:

```bash
go build
```

## License

LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
