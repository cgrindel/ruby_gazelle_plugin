#!/usr/bin/env bash

# --- begin runfiles.bash initialization v3 ---
# Copy-pasted from the Bazel Bash runfiles library v3.
set -uo pipefail
set +e
f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null \
  || source "$(grep -sm1 "^$f " "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" 2>/dev/null \
  || source "$0.runfiles/$f" 2>/dev/null \
  || source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null \
  || source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null \
  || {
    echo >&2 "ERROR: cannot find $f"
    exit 1
  }
f=
set -e
# --- end runfiles.bash initialization v3 ---

# MARK - Deps

yamlfmt_location=rules_multitool++multitool+multitool/tools/yamlfmt/yamlfmt
yamlfmt="$(rlocation "${yamlfmt_location}")" \
  || (echo >&2 "Failed to locate ${yamlfmt_location}" && exit 1)

yamlfmt_yml_location=_main/.yamlfmt.yml
yamlfmt_yml="$(rlocation "${yamlfmt_yml_location}")" \
  || (echo >&2 "Failed to locate ${yamlfmt_yml_location}" && exit 1)

# MARK - Execute command

cmd=(
  "${yamlfmt}"
  -conf "${yamlfmt_yml}"
)
if [[ $# -gt 0 ]]; then
  cmd+=("$@")
fi
"${cmd[@]}"
