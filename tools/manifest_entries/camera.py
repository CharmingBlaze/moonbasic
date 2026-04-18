"""Camera manifest entries.

Covers: CAMERA.*, CAMERA2D.*, camera Blitz-compat globals.
Arities verified from runtime/camera source.
"""

from .helpers import ntimes

ENTRIES = [
    ("CAMERA.FREE", ntimes("handle", 1), None, "rendering"),
    ("CAMERA.GETPOS", ntimes("handle", 1), "handle", "rendering"),
    ("CAMERA.GETTARGET", ntimes("handle", 1), "handle", "rendering"),
    ("CAMERA.SETUP", ["handle", "float", "float", "float"], None, "rendering"),
    ("CAMERA2D.GETMATRIX", ntimes("handle", 1), "handle", "rendering"),
    ("CAMERA2D.GETOFFSET", ntimes("handle", 1), "array", "rendering"),
    ("CAMERA2D.GETZOOM", ntimes("handle", 1), "float", "rendering"),
    ("CAMERA2D.SCREENTOWORLD", ["handle", "float", "float"], "handle", "rendering"),
    ("CAMERA2D.WORLDTOSCREEN", ["handle", "float", "float"], "handle", "rendering"),
    # --- camera extended (Blitz-compat) ---
    ("CAMERAFOGMODE", ["handle", "int"], None, "rendering"),
    ("CAMERAFOGCOLOR", ["handle", "int", "int", "int"], None, "rendering"),
    ("CAMERAFOGRANGE", ["handle", "float", "float"], None, "rendering"),
    ("CAMERAPROJECT", ["handle", "float", "float", "float"], None, "rendering"),
    ("CAMERARANGE", ["handle", "float", "float"], None, "rendering"),
    ("CAMERAVIEWPORT", ["handle", "int", "int", "int", "int"], None, "rendering"),
    ("MOVECAMERA", ["handle", "float", "float", "float"], None, "rendering"),
    ("TURNCAMERA", ["handle", "float", "float", "float"], None, "rendering"),
    ("SHAKECAMERA", ["handle", "float", "float"], None, "rendering"),
]
