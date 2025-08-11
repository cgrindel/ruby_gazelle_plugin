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

shfmt_location=rules_multitool++multitool+multitool/tools/shfmt/shfmt
shfmt="$(rlocation "${shfmt_location}")" \
  || (echo >&2 "Failed to locate ${shfmt_location}" && exit 1)

# MARK - Execute command

cmd=(
  "${shfmt}"
  --indent 2
  --case-indent
  --binary-next-line
)
if [[ $# -gt 0 ]]; then
  cmd+=("$@")
fi
"${cmd[@]}"
