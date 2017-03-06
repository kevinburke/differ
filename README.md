# differ

Differ makes it easy to run a command and error if it generated a change in a
git worktree. You can use this in tests or the build process to verify that
a given build step was run correctly. For example you may want to verify that
all files in a Go project have run `go fmt`. Run:

```
differ go fmt ./...
```

This will execute `go fmt ./...` and error if it modifies any file tracked by
Git.

Other uses:

- Restore and revendor all vendored libraries and error if a git diff is
generated.
- Check whether new CSS files have been generated from SCSS, HTML files from
  Markdown, JS files from Coffeescript, or any other compilation step.

## Usage

Run the same command you would usually run but put `differ` before it. differ
will exit with a non-zero return code if:

- your command exits with an error

- "git status" errors

- "git status" says that there are untracked or modified files present

## Installation
