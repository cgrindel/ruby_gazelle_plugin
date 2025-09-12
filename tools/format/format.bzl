"""Implementation for `format` macro."""

load("@aspect_rules_lint//format:defs.bzl", "format_multirun", "format_test")

def format(name, **kwargs):
    """Defines `format_multirun` and `format_test` targets based upon the \
    provided attributes.

    Note: The generated `format_test` target runs in non-hermetic mode (no_sandbox)
    to format the entire workspace. This will be replaced with hermetic tests once
    a Gazelle plugin is available.

    Args:
        name: The name for that `format_multirun` target as a `string`. This
            value is the base name for the `format_test`.
        **kwargs: The common attributes to pass to the generated targets.
    """
    format_multirun(
        name = name,
        **kwargs
    )

    # NOTE: This is a non-hermetic test. It runs outside the sandbox accessing
    # places that normal tests should never access. We will replace this with
    # hermetic tests, once a Gazelle plugin is ready to do the heavy lifting.
    # https://github.com/aspect-build/rules_lint/blob/main/docs/formatting.md#2-test-target
    format_test(
        name = "{name}_test".format(name = name),
        # Enables formatting the entire workspace, paired with 'workspace'
        # attribute
        no_sandbox = True,
        # A file in the workspace root, where the no_sandbox mode will run the
        # formatter
        workspace = "//:MODULE.bazel",
        **kwargs
    )
