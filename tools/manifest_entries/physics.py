"""Physics and character-controller manifest entries.

Covers: PHYSICS collision info globals, PHYSICS.TORQUE, PHYSICS2D.*,
        CHARACTERREF.*.
Arities verified from runtime/physics3d and runtime/mbentity source.
"""

from .helpers import ntimes

ENTRIES = [
    # --- physics collision info ---
    ("PHYSICSCOLLISIONFORCE", ["handle", "int"], "float", "any"),
    ("PHYSICSCOLLISIONPX", ["handle", "int"], "float", "any"),
    ("PHYSICSCOLLISIONPY", ["handle", "int"], "float", "any"),
    ("PHYSICSCOLLISIONPZ", ["handle", "int"], "float", "any"),
    ("PHYSICSCOLLISIONNX", ["handle", "int"], "float", "any"),
    ("PHYSICSCOLLISIONNY", ["handle", "int"], "float", "any"),
    ("PHYSICSCOLLISIONNZ", ["handle", "int"], "float", "any"),
    ("PHYSICSCOLLISIONY", ["handle", "int"], "float", "any"),
    ("PHYSICSCONTACTCOUNT", ntimes("handle", 1), "int", "any"),
    ("COLLISIONFORCE", ["handle", "int"], "float", "any"),
    ("PHYSICS.TORQUE", ["handle", "float", "float", "float"], None, "any"),
    # --- physics2d extended ---
    ("PHYSICS2D.ONCOLLISION", ["handle", "string"], None, "any"),
    ("PHYSICS2D.PROCESSCOLLISIONS", [], None, "any"),
    # --- characterref ---
    ("CHARACTERREF.GETGROUNDSTATE", ntimes("handle", 1), "int", "any"),
    ("CHARACTERREF.GETVELOCITY", ntimes("handle", 1), "handle", "any"),
    ("CHARACTERREF.ISMOVING", ntimes("handle", 1), "bool", "any"),
    ("CHARACTERREF.MOVEWITHCAMERA", ["handle", "float", "float", "handle"], None, "any"),
    ("CHARACTERREF.SETBOUNCE", ["handle", "float"], None, "any"),
    ("CHARACTERREF.SETGRAVITYSCALE", ["handle", "float"], None, "any"),
    ("CHARACTERREF.SETLINEARVELOCITY", ["handle", "float", "float", "float"], None, "any"),
    ("CHARACTERREF.SETPADDING", ["handle", "float"], None, "any"),
    ("CHARACTERREF.SETSTICKDOWN", ["handle", "bool"], None, "any"),
    ("CHARACTERREF.UPDATEMOVE", ntimes("handle", 1), None, "any"),
]
