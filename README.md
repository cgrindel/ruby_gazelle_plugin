# Ruby Gazelle Plugin

A Bazel Gazelle plugin for Ruby that automatically generates `rb_library`, `rb_binary`, and
`rb_test` targets from Ruby source files.

## Overview

This plugin extends [Gazelle](https://github.com/bazelbuild/bazel-gazelle) to support Ruby
projects in Bazel. It scans Ruby source files and generates appropriate build targets, making it
easier to integrate Ruby projects with Bazel.

## Features

- Automatic generation of Ruby build targets
- Support for `rb_library`, `rb_binary`, and `rb_test` rules
- Integration with existing Bazel and Gazelle workflows

## Installation

Add the following to your `MODULE.bazel`:

```starlark
bazel_dep(name = "ruby_gazelle_plugin", version = "0.0.0")
```

## Usage

To use the Ruby Gazelle plugin:

1. Add the plugin to your `gazelle_binary` in your root `BUILD.bazel`:

```starlark
load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle_bin",
    languages = [
        "@rules_go//go/tools/builders:default_go_tool",
        "@ruby_gazelle_plugin//gazelle",
    ],
)

gazelle(
    name = "update_build_files",
    gazelle = ":gazelle_bin",
)
```

2. Run Gazelle to generate BUILD files:

```bash
bazel run //:update_build_files
```

## Requirements

- Bazel with bzlmod support
- [rules_ruby](https://github.com/bazelruby/rules_ruby) for Ruby rule definitions

## Development

This plugin is implemented in Go and follows the standard Gazelle plugin architecture.
