"""Model, material, mesh, and animation manifest entries.

Covers: MODEL.LOADASYNC, MODEL.ISLOADED, MODEL.MAKEBOX, MODEL.MAKECAPSULE,
        MATERIAL.*, ANIMLENGTH, MESH.* (creation entries already in main
        manifest; these are the missing extras).
Arities verified from runtime/mbmodel3d source.
"""

from .helpers import ntimes

ENTRIES = [
    ("MODEL.LOADASYNC", ntimes("string", 1), "handle", "any"),
    ("MODEL.ISLOADED", ntimes("handle", 1), "bool", "any"),
    ("MODEL.MAKEBOX", ntimes("float", 3), "handle", "any"),
    ("MODEL.MAKECAPSULE", ntimes("float", 2), "handle", "any"),
    ("MATERIAL.CREATEDEFAULT", [], "handle", "any"),
    ("MATERIAL.CREATEPBR", [], "handle", "any"),
    ("ANIMLENGTH", ntimes("handle", 1), "int", "any"),
]
