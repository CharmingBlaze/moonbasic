# Gesture Commands

Touch gesture detection for mobile/tablet: tap, double-tap, hold, swipe, drag, and pinch. Works alongside `INPUT.*` touch commands.

## Core Workflow

1. `GESTURE.ENABLE(gestureMask)` — enable the gestures you need (bitmask).
2. Each frame: `GESTURE.GETDETECTED()` → compare with gesture constants.
3. Read gesture data with `GESTURE.GETHOLDDURATION`, `GESTURE.GETDRAGVECTORX/Y`, `GESTURE.GETPINCHVECTORX/Y/ANGLE`.

Gesture constants (bitmask): `GESTURE_TAP`=1, `GESTURE_DOUBLETAP`=2, `GESTURE_HOLD`=4, `GESTURE_DRAG`=8, `GESTURE_SWIPE_RIGHT`=16, `GESTURE_SWIPE_LEFT`=32, `GESTURE_SWIPE_UP`=64, `GESTURE_SWIPE_DOWN`=128, `GESTURE_PINCH_IN`=256, `GESTURE_PINCH_OUT`=512.

---

## Enable

### `GESTURE.ENABLE(mask)` 

Enables gesture recognition for the given bitmask. Pass `-1` or `1023` to enable all.

---

## Detection

### `GESTURE.GETDETECTED()` 

Returns the bitmask of gestures detected this frame. Compare with gesture constants.

---

### `GESTURE.ISDETECTED(gesture)` 

Returns `TRUE` if the given single gesture constant was detected this frame.

---

## Hold

### `GESTURE.GETHOLDDURATION()` 

Returns how long (seconds) the hold gesture has been active this frame.

---

## Drag

### `GESTURE.GETDRAGVECTORX()` / `GESTURE.GETDRAGVECTORY()` 

Returns the drag vector components for the current drag gesture.

---

### `GESTURE.GETDRAGANGLE()` 

Returns the drag direction in degrees.

---

## Pinch

### `GESTURE.GETPINCHVECTORX()` / `GESTURE.GETPINCHVECTORY()` 

Returns the pinch vector components.

---

### `GESTURE.GETPINCHANGLE()` 

Returns the pinch rotation angle in degrees.

---

## Full Example

Detecting tap and swipe gestures on a mobile target.

```basic
WINDOW.OPEN(480, 854, "Gesture Demo")
WINDOW.SETFPS(60)

GESTURE.ENABLE(GESTURE_TAP + GESTURE_SWIPE_LEFT + GESTURE_SWIPE_RIGHT + GESTURE_HOLD)

msg = "Touch the screen"

WHILE NOT WINDOW.SHOULDCLOSE()
    g = GESTURE.GETDETECTED()

    IF GESTURE.ISDETECTED(GESTURE_TAP) THEN
        msg = "TAP!"
    END IF
    IF GESTURE.ISDETECTED(GESTURE_SWIPE_LEFT) THEN
        msg = "Swipe LEFT"
    END IF
    IF GESTURE.ISDETECTED(GESTURE_SWIPE_RIGHT) THEN
        msg = "Swipe RIGHT"
    END IF
    IF GESTURE.ISDETECTED(GESTURE_HOLD) THEN
        msg = "Hold: " + STR(GESTURE.GETHOLDDURATION()) + "s"
    END IF

    RENDER.CLEAR(20, 20, 40)
    DRAW.TEXT(msg, 80, 400, 32, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [INPUT.md](INPUT.md) — touch point positions and count
- [MOUSE.md](MOUSE.md) — mouse button and position (desktop)
