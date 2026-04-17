"""Shared helpers for manifest entry definitions."""


def ntimes(k: str, n: int) -> list[str]:
    """Return a list of *n* copies of type string *k*."""
    return [k] * n
