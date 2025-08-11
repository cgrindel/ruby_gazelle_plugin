# Contributing to `ruby_gazelle_plugin`

This guide will help you get set up and understand our development process.

## Table of Contents

- [Development Workflow](#development-workflow)
- [Pull Request Procedure](#pull-request-procedure)
- [Continuous Integration](#continuous-integration)
- [Coding Styles](#coding-styles)
  - [Bazel Starlark](#bazel-starlark)
  - [Go](#go)
  - [Markdown](#markdown)
  - [Shell](#shell)
  - [YAML](#yaml)
- [Optional Tools](#optional-tools)

## Development Workflow

After completing the setup, familiarize yourself with these essential commands:

```bash
# Format and update all source files (comprehensive)
bazel run //:tidy

# Format source files only
bazel run //:format

# Update source files without formatting
bazel run //:update_files

# Check if formatting is needed (used in CI)
bazel run //:tidy_check

# Run all tests
bazel test //...

# Build all targets
bazel build //...
```

For more details about the project structure and development practices,
see the [CLAUDE.md](CLAUDE.md) file.

## Pull Request Procedure

- Create new pull requests as a draft assigned to yourself.
- Pull request titles must follow
  [conventional commit syntax][conventional-commits].
- Related GitHub issues should be [linked to the pull request][gh-linking].
- Mark the pull request as Ready to Review and request a review once you have
  seen all CI tests succeed.
- Pull requests must be reviewed by a codeowner before merging to the target
  branch.

## Continuous Integration

This repository uses [GitHub Actions][github-actions] as
its CI platform. The CI configuration file is located at
`.github/workflows/ci.yml`. CI is configured to run on all pull requests
that are targeted to be merged to `main`.

## Coding Styles

### Bazel Starlark

This repository follows [Google's Starlark style guide][starlark-style].
Use buildifier to format and lint Starlark files.

### Go

This repository follows [Google's Go style guide][go-style].
The maximum line length for Go files should be 80 characters.
Use `gofmt` or `goimports` for formatting.

### Markdown

We follow [Google's Markdown style guide][markdown-style].
We use [prettier][prettier] for formatting Markdown files.

### Shell

We follow [Google's Shell style guide][shell-style].
We use [shfmt][shfmt] to format and
[shellcheck][shellcheck] for linting.

### YAML

We use [yamlfmt][yamlfmt] to format YAML files.
There is no formal style guide.

## Optional Tools

These tools are not required but can enhance your development experience:

### Claude Code

[Claude Code][claude-code] is an AI-powered CLI
tool for software engineering tasks. It can help with code review, refactoring,
debugging, and understanding complex codebases.

Claude Code is particularly useful for:

- Analyzing and explaining code structure
- Automated code reviews and suggestions
- Generating documentation and tests
- Debugging complex issues
- Learning new codebases quickly

This repository includes pre-configured settings for Claude Code:

- **CLAUDE.md**: Project-specific guidance and development practices
- **.claude/**: Additional configuration files (if present)

[claude-code]: https://www.anthropic.com/claude-code
[conventional-commits]: https://www.conventionalcommits.org/en/v1.0.0/
[gh-linking]: https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue
[git-worktree-docs]: https://git-scm.com/docs/git-worktree
[git-worktree-tutorial]: https://www.google.com/search?q=git+worktree+tutorial
[github-actions]: https://github.com/features/actions
[markdown-style]: https://google.github.io/styleguide/docguide/style.html
[prettier]: https://prettier.io/
[go-style]: https://google.github.io/styleguide/go/
[shell-style]: https://google.github.io/styleguide/shellguide.html
[shellcheck]: https://www.shellcheck.net/
[shfmt]: https://github.com/mvdan/sh
[starlark-style]: https://bazel.build/rules/bzl-style
[yamlfmt]: https://github.com/google/yamlfmt
