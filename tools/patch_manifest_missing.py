#!/usr/bin/env python3
"""Append runtime-registered commands missing from compiler/builtinmanifest/commands.json."""
import json
import os

ROOT = os.path.join(os.path.dirname(__file__), "..")
MAN = os.path.join(ROOT, "compiler", "builtinmanifest", "commands.json")


def ntimes(k, n):
    return [k] * n


# (key, args list, optional returns string, optional phase)
# Arities from runtime handlers (expects / len(args) checks).
NEW = [
    # --- audio ---
    ("AUDIO.GETMUSICLENGTH", ntimes("handle", 1), "float", "any"),
    ("AUDIO.GETMUSICTIME", ntimes("handle", 1), "float", "any"),
    ("AUDIO.ISMUSICPLAYING", ntimes("handle", 1), "bool", "any"),
    ("AUDIO.ISSOUNDPLAYING", ntimes("handle", 1), "bool", "any"),
    ("AUDIO.SEEKMUSIC", ["handle", "float"], None, "any"),
    ("AUDIO.SETMASTERVOLUME", ntimes("float", 1), None, "any"),
    ("AUDIO.SETMUSICPITCH", ["handle", "float"], None, "any"),
    ("AUDIO.SETMUSICVOLUME", ["handle", "float"], None, "any"),
    ("AUDIO.SETSOUNDPAN", ["handle", "float"], None, "any"),
    ("AUDIO.SETSOUNDPITCH", ["handle", "float"], None, "any"),
    ("AUDIO.SETSOUNDVOLUME", ["handle", "float"], None, "any"),
    ("AUDIO.UPDATEMUSIC", ntimes("handle", 1), None, "any"),
    # --- camera ---
    ("CAMERA.FREE", ntimes("handle", 1), None, "rendering"),
    ("CAMERA.GETPOS", ntimes("handle", 1), "handle", "rendering"),
    ("CAMERA.GETTARGET", ntimes("handle", 1), "handle", "rendering"),
    ("CAMERA.SETUP", ["handle", "float", "float", "float"], None, "rendering"),
    ("CAMERA2D.GETMATRIX", ntimes("handle", 1), "handle", "rendering"),
    ("CAMERA2D.SCREENTOWORLD", ["handle", "float", "float"], "handle", "rendering"),
    ("CAMERA2D.WORLDTOSCREEN", ["handle", "float", "float"], "handle", "rendering"),
    # --- draw 2d (pixel-ish coords: int; floats where runtime uses argFloat) ---
    ("DRAW.ARC", ntimes("float", 6) + ntimes("int", 4), None, "rendering"),
    (
        "DRAW.BILLBOARD",
        ["handle", "float", "float", "float", "float"] + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.BILLBOARDREC",
        ["handle", "float", "float", "float", "float", "float", "float", "float", "float"]
        + ntimes("int", 4),
        None,
        "rendering",
    ),
    ("DRAW.BOUNDINGBOX", ntimes("float", 6) + ntimes("int", 4), None, "rendering"),
    (
        "DRAW.CAPSULE",
        ntimes("float", 9) + ntimes("int", 2) + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.CAPSULEWIRES",
        ntimes("float", 9) + ntimes("int", 2) + ntimes("int", 4),
        None,
        "rendering",
    ),
    ("DRAW.CIRCLE", ntimes("int", 7), None, "rendering"),
    ("DRAW.CIRCLEGRADIENT", ntimes("int", 3) + ntimes("int", 8), None, "rendering"),
    ("DRAW.CIRCLELINES", ntimes("int", 7), None, "rendering"),
    ("DRAW.CIRCLESECTOR", ntimes("int", 6) + ntimes("int", 4), None, "rendering"),
    ("DRAW.CUBE", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW.CUBEWIRES", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW.CYLINDER", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW.CYLINDERWIRES", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW.DOT", ntimes("int", 2) + ntimes("float", 1) + ntimes("int", 4), None, "rendering"),
    ("DRAW.ELLIPSE", ntimes("int", 8), None, "rendering"),
    ("DRAW.ELLIPSELINES", ntimes("int", 8), None, "rendering"),
    ("DRAW.GETPIXELCOLOR", ntimes("int", 2), "handle", "rendering"),
    ("DRAW.GRID", ["int", "float"], None, "rendering"),
    ("DRAW.GRID2D", ntimes("int", 5), None, "rendering"),
    ("DRAW.LINE", ntimes("int", 8), None, "rendering"),
    ("DRAW.LINE3D", ntimes("float", 6) + ntimes("int", 4), None, "rendering"),
    ("DRAW.LINEBEZIER", ntimes("float", 5) + ntimes("int", 4), None, "rendering"),
    ("DRAW.LINEBEZIERCUBIC", ntimes("float", 9) + ntimes("int", 4), None, "rendering"),
    ("DRAW.LINEBEZIERQUAD", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW.LINEEX", ntimes("float", 5) + ntimes("int", 4), None, "rendering"),
    ("DRAW.PIXEL", ntimes("int", 6), None, "rendering"),
    ("DRAW.PIXELV", ntimes("float", 2) + ntimes("int", 4), None, "rendering"),
    ("DRAW.PLANE", ntimes("float", 6) + ntimes("int", 4), None, "rendering"),
    ("DRAW.POINT3D", ntimes("float", 3) + ntimes("int", 4), None, "rendering"),
    ("DRAW.POLY", ntimes("int", 5) + ntimes("int", 4), None, "rendering"),
    ("DRAW.POLYLINES", ntimes("int", 5) + ntimes("float", 1) + ntimes("int", 4), None, "rendering"),
    ("DRAW.RAY", ["handle"] + ntimes("int", 4), None, "rendering"),
    ("DRAW.RECTGRAD", ntimes("int", 4) + ntimes("int", 16), None, "rendering"),
    ("DRAW.RECTGRADH", ntimes("int", 12), None, "rendering"),
    ("DRAW.RECTGRADV", ntimes("int", 12), None, "rendering"),
    ("DRAW.RECTLINES", ntimes("int", 5) + ntimes("float", 1) + ntimes("int", 4), None, "rendering"),
    ("DRAW.RECTPRO", ntimes("int", 4) + ntimes("float", 3) + ntimes("int", 4), None, "rendering"),
    ("DRAW.RING", ntimes("int", 7) + ntimes("float", 2) + ntimes("int", 4), None, "rendering"),
    ("DRAW.RINGLINES", ntimes("int", 7) + ntimes("float", 2) + ntimes("int", 4), None, "rendering"),
    ("DRAW.SETPIXELCOLOR", ntimes("int", 6), None, "rendering"),
    ("DRAW.SPHERE", ntimes("float", 4) + ntimes("int", 4), None, "rendering"),
    ("DRAW.SPHEREWIRES", ntimes("float", 4) + ntimes("int", 2) + ntimes("int", 4), None, "rendering"),
    ("DRAW.SPLINEBASIS", ["handle", "float"] + ntimes("int", 4), None, "rendering"),
    ("DRAW.SPLINEBEZIERCUBIC", ["handle", "float"] + ntimes("int", 4), None, "rendering"),
    ("DRAW.SPLINEBEZIERQUAD", ["handle", "float"] + ntimes("int", 4), None, "rendering"),
    ("DRAW.SPLINECATMULLROM", ["handle", "float"] + ntimes("int", 4), None, "rendering"),
    ("DRAW.SPLINELINEAR", ["handle", "float"] + ntimes("int", 4), None, "rendering"),
    ("DRAW.TEXT", ["string"] + ntimes("int", 7), None, "rendering"),
    (
        "DRAW.TEXTEX",
        ["handle", "string", "float", "float", "float", "float"] + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.TEXTFONT",
        ["handle", "string", "float", "float", "float", "float"] + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.TEXTFONTWIDTH",
        ["handle", "string", "float", "float"],
        "float",
        "rendering",
    ),
    (
        "DRAW.TEXTPRO",
        ["handle", "string", "float", "float", "float", "float", "float", "float", "float"]
        + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.TEXTUREEX",
        ["handle", "float", "float", "float", "float"] + ntimes("int", 4),
        None,
        "rendering",
    ),
    ("DRAW.TEXTUREFLIPPED", ntimes("handle", 1), None, "rendering"),
    ("DRAW.TEXTUREFULL", ntimes("handle", 1), None, "rendering"),
    (
        "DRAW.TEXTUREPRO",
        ["handle", "float", "float", "float", "float", "float", "float", "float", "float", "float", "float"]
        + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.TEXTUREREC",
        ["handle", "float", "float", "float", "float", "float", "float"] + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.TEXTURETILED",
        ["handle"] + ntimes("float", 12) + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW.TEXTUREV",
        ["handle", "float", "float"] + ntimes("int", 4),
        None,
        "rendering",
    ),
    ("DRAW.TEXTWIDTH", ["string", "int"], "int", "rendering"),
    ("DRAW.TRIANGLE", ntimes("int", 6) + ntimes("int", 4), None, "rendering"),
    ("DRAW.TRIANGLELINES", ntimes("int", 6) + ntimes("int", 4), None, "rendering"),
    # --- draw 3d ---
    ("DRAW3D.BBOX", ntimes("float", 6) + ntimes("int", 4), None, "rendering"),
    (
        "DRAW3D.BILLBOARD",
        ["handle", "float", "float", "float", "float"] + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW3D.BILLBOARDREC",
        ["handle", "float", "float", "float", "float", "float", "float", "float", "float"]
        + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW3D.CAPSULE",
        ntimes("float", 9) + ntimes("int", 2) + ntimes("int", 4),
        None,
        "rendering",
    ),
    (
        "DRAW3D.CAPSULEWIRES",
        ntimes("float", 9) + ntimes("int", 2) + ntimes("int", 4),
        None,
        "rendering",
    ),
    ("DRAW3D.CUBE", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.CUBEWIRES", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.CYLINDER", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.CYLINDERWIRES", ntimes("float", 7) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.GRID", ["int", "float"], None, "rendering"),
    ("DRAW3D.LINE", ntimes("float", 6) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.PLANE", ntimes("float", 6) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.POINT", ntimes("float", 3) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.RAY", ["handle"] + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.SPHERE", ntimes("float", 4) + ntimes("int", 4), None, "rendering"),
    ("DRAW3D.SPHEREWIRES", ntimes("float", 4) + ntimes("int", 2) + ntimes("int", 4), None, "rendering"),
    # --- text helpers (also reachable as bare calls; manifest for docs / tooling) ---
    ("GETTEXTCODEPOINTCOUNT", ntimes("string", 1), "int", "any"),
    ("MEASURETEXT", ["string", "int"], "int", "any"),
    ("MEASURETEXTEX", ["handle", "string", "float", "float"], "handle", "any"),
    # --- input ---
    ("INPUT.CHARPRESSED", [], "int", "any"),
    ("INPUT.GETGAMEPADAXISVALUE", ntimes("int", 2), "float", "any"),
    ("INPUT.ISGAMEPADAVAILABLE", ntimes("int", 1), "bool", "any"),
    ("INPUT.MOUSEDELTAX", [], "float", "any"),
    ("INPUT.MOUSEDELTAY", [], "float", "any"),
    ("INPUT.MOUSEPRESSED", ntimes("int", 1), "bool", "any"),
    ("INPUT.MOUSERELEASED", ntimes("int", 1), "bool", "any"),
    ("INPUT.MOUSEWHEELMOVE", [], "float", "any"),
    ("INPUT.SETMOUSEPOS", ntimes("int", 2), None, "any"),
    # --- misc ---
    ("MUSIC.FREE", ntimes("handle", 1), None, "any"),
    ("RENDER.SET2DAMBIENT", ntimes("int", 4), None, "rendering"),
    # --- texture ---
    (
        "TEXTURE.GENCHECKED",
        ntimes("int", 4) + ntimes("handle", 2),
        "handle",
        "any",
    ),
    ("TEXTURE.GENCOLOR", ntimes("int", 6), "handle", "any"),
    ("TEXTURE.GENGRADIENTH", ntimes("int", 2) + ntimes("handle", 2), "handle", "any"),
    ("TEXTURE.GENGRADIENTV", ntimes("int", 2) + ntimes("handle", 2), "handle", "any"),
    ("TEXTURE.HEIGHT", ntimes("handle", 1), "int", "any"),
    ("TEXTURE.SETFILTER", ["handle", "int"], None, "any"),
    ("TEXTURE.SETWRAP", ["handle", "int"], None, "any"),
    ("TEXTURE.UPDATE", ntimes("handle", 2), None, "any"),
    ("TEXTURE.WIDTH", ntimes("handle", 1), "int", "any"),
    ("TIME.GETFPS", [], "float", "any"),
    # --- util ---
    ("UTIL.COPYFILE", ntimes("string", 2), None, "any"),
    ("UTIL.DELETEDIR", ntimes("string", 1), None, "any"),
    ("UTIL.DELETEFILE", ntimes("string", 1), None, "any"),
    ("UTIL.GETDIR", [], "string", "any"),
    ("UTIL.GETDIRS", ntimes("string", 1), "string", "any"),
    ("UTIL.MOVEFILE", ntimes("string", 2), None, "any"),
    ("UTIL.RENAMEFILE", ntimes("string", 2), None, "any"),
    # --- window ---
    ("WINDOW.GETFPS", [], "int", "any"),
    ("WINDOW.HEIGHT", [], "int", "any"),
    ("WINDOW.ISFULLSCREEN", [], "bool", "any"),
    ("WINDOW.ISRESIZED", [], "bool", "any"),
    ("WINDOW.TOGGLEFULLSCREEN", [], None, "any"),
    ("WINDOW.WIDTH", [], "int", "any"),
]


def main():
    with open(MAN, encoding="utf-8") as f:
        root = json.load(f)
    have = {c["key"] for c in root["commands"]}
    added = 0
    for spec in NEW:
        key = spec[0]
        if key in have:
            continue
        args = spec[1]
        rec = {"key": key, "args": args, "phase": spec[3]}
        if len(spec) > 2 and spec[2]:
            rec["returns"] = spec[2]
        root["commands"].append(rec)
        added += 1
        have.add(key)
    with open(MAN, "w", encoding="utf-8", newline="\n") as f:
        json.dump(root, f, indent=2)
        f.write("\n")
    print("added", added, "commands")


if __name__ == "__main__":
    main()
