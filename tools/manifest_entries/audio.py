"""Audio, music, sound, and font-loading manifest entries.

Covers: AUDIO.*, MUSIC.*, LOADMUSIC, PLAYMUSIC, STOPMUSIC, LOOPSOUND,
        STOPSOUND, SOUNDPAN, LOADFONT, FONT.*, LOADIMAGE.
Arities verified from runtime/audio and runtime/mbfont source.
"""

from .helpers import ntimes

ENTRIES = [
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
    ("AUDIO.PLAYRNDSOUND", ["handle", "int"], None, "any"),
    ("AUDIO.PLAYVARYSOUND", ["handle", "float", "float"], None, "any"),
    # --- audio extended (Blitz-style globals) ---
    ("LOADMUSIC", ntimes("string", 1), "handle", "any"),
    ("PLAYMUSIC", ntimes("handle", 1), None, "any"),
    ("STOPMUSIC", ntimes("handle", 1), None, "any"),
    ("MUSICVOLUME", ["handle", "float"], None, "any"),
    ("LOOPSOUND", ntimes("handle", 1), None, "any"),
    ("STOPSOUND", ntimes("handle", 1), None, "any"),
    ("SOUNDPAN", ["handle", "float"], None, "any"),
    ("MUSIC.FREE", ntimes("handle", 1), None, "any"),
    # --- font / image loading ---
    ("LOADFONT", ["string", "int"], "handle", "any"),
    ("FONT.SETDEFAULT", ntimes("handle", 1), None, "any"),
    ("LOADIMAGE", ntimes("string", 1), "handle", "any"),
]
