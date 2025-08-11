## Markdown Guidelines

- The maximum line length for markdown files is 80 characters.

## Style Guide

- We will use Google's Go style guide for this repository. It is located at
  https://google.github.io/styleguide/go/.
- The maximum line length for Go files should be 80 characters.

## Development Workflow

- See [CONTRIBUTING.md](CONTRIBUTING.md) for comprehensive development workflow,
  formatting commands, and complete style guides for all languages used in this
  project.

- Confirm the build works properly by running bazel test //...

- When we add or remove Go source files, be sure to update the corresponding
  BUILD.bazel file.
- When make changes to the import statement for a Go source file, be sure to
  update the corresponding BUILD.bazel file.
- After making changes to any source file, run `bazel test //...` to ensure
  everything works properly.
