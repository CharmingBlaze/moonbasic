# Color helpers (`RGB`, `ARGB`, mix)

Registered from **`runtime/mbgame`** (see **`register_color_format.go`**, **`pure_color.go`**). Use these for **integer** color construction and formatting — drawing still uses **`r, g, b, a` byte tuples** on **`Draw.*`** in many places.

Common entries:

- **`RGB(r, g, b)`** / **`ARGB(a, r, g, b)`** — packed integers.
- **`RGBMIX`**, grayscale, and related helpers — see **`runtime/mbgame`** registrations.

Formatting text for HUDs (**`FORMATINT`**, **`FORMATTIME`**, …) is also on **`mbgame`**; prefer calling those rather than duplicating format logic in draw helpers.
