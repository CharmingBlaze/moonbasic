"""Manifest entry modules for patch_manifest_missing.py.

Each submodule exports an ENTRIES list of (key, args, returns, phase) tuples
representing runtime-registered commands that need compiler manifest entries.
Arg types and counts are verified from runtime Go source (len(args) checks).

Modules are split by namespace group to keep files focused and reviewable.
See STYLE_GUIDE.md §Documentation and §Naming Conventions.
"""

from .helpers import ntimes  # noqa: F401 — re-export for submodules
