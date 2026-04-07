# Complete MoonBasic Command Documentation

This document provides comprehensive documentation for ALL MoonBasic commands with their exact arguments and expected types.

## Audio Commands (mbaudio)

### AUDIO.LOADMUSIC
**Arguments:** `path$`
- `path$` (string): Path to music file

### AUDIO.UPDATEMUSIC
**Arguments:** `musicHandle`
- `musicHandle` (handle): Music handle to update

### MUSIC.FREE
**Arguments:** `musicHandle`
- `musicHandle` (handle): Music handle to free

### AUDIO.LOADSOUND
**Arguments:** `path$`
- `path$` (string): Path to sound file

### FREESOUND
**Arguments:** `soundHandle`
- `soundHandle` (handle): Sound handle to free

### LOADSOUND (Easy Mode)
**Arguments:** `path$`
- `path$` (string): Path to sound file

## Window Commands

### WINDOW.OPEN
**Arguments:** `width, height, title$`
- `width` (numeric): Window width in pixels
- `height` (numeric): Window height in pixels  
- `title$` (string): Window title

### WINDOW.CANOPEN
**Arguments:** `width, height, title$`
- `width` (numeric): Window width in pixels
- `height` (numeric): Window height in pixels
- `title$` (string): Window title

### WINDOW.SETFPS
**Arguments:** `fps`
- `fps` (numeric): Target frames per second

### WINDOW.CLOSE
**Arguments:** None

### WINDOW.SHOULDCLOSE
**Arguments:** None

### WINDOW.SETFLAG
**Arguments:** `flag`
- `flag` (int): Window flag (use FLAG_* constants)

### WINDOW.CLEARFLAG
**Arguments:** `flag`
- `flag` (int): Window flag to clear (use FLAG_* constants)

### WINDOW.CHECKFLAG
**Arguments:** `flag`
- `flag` (int): Window flag to check (use FLAG_* constants)

### WINDOW.SETSTATE
**Arguments:** `flag`
- `flag` (int): Window state flag (use FLAG_* constants)

### WINDOW.SETMINSIZE
**Arguments:** `width, height`
- `width` (numeric): Minimum window width
- `height` (numeric): Minimum window height

### WINDOW.SETMAXSIZE
**Arguments:** `width, height`
- `width` (numeric): Maximum window width
- `height` (numeric): Maximum window height

### WINDOW.GETPOSITIONX
**Arguments:** None

### WINDOW.GETPOSITIONY
**Arguments:** None

### WINDOW.SETMONITOR
**Arguments:** `monitorIndex`
- `monitorIndex` (int): Monitor index to use

### WINDOW.GETMONITORCOUNT
**Arguments:** None

### WINDOW.GETMONITORNAME
**Arguments:** `monitorIndex`
- `monitorIndex` (int): Monitor index

### WINDOW.GETMONITORWIDTH
**Arguments:** `monitorIndex`
- `monitorIndex` (int): Monitor index

### WINDOW.GETMONITORHEIGHT
**Arguments:** `monitorIndex`
- `monitorIndex` (int): Monitor index

### WINDOW.GETMONITORREFRESHRATE
**Arguments:** `monitorIndex`
- `monitorIndex` (int): Monitor index

### WINDOW.GETSCALEDPIX
**Arguments:** None

### WINDOW.GETSCALEDPIY
**Arguments:** None

### WINDOW.SETICON
**Arguments:** `imageHandle`
- `imageHandle` (handle): Image handle for window icon

### WINDOW.SETOPACITY
**Arguments:** `opacity#`
- `opacity#` (float): Window opacity (0.0-1.0)

### WINDOW.SETPOSITION
**Arguments:** `x, y`
- `x` (numeric): Window X position
- `y` (numeric): Window Y position

### WINDOW.SETSIZE
**Arguments:** `width, height`
- `width` (numeric): Window width
- `height` (numeric): Window height

### WINDOW.MINIMIZE
**Arguments:** None

### WINDOW.MAXIMIZE
**Arguments:** None

### WINDOW.RESTORE
**Arguments:** None

### WINDOW.SETTARGETFPS
**Arguments:** `fps`
- `fps` (numeric): Target frames per second

### WINDOW.WIDTH
**Arguments:** None

### WINDOW.HEIGHT
**Arguments:** None

### WINDOW.GETFPS
**Arguments:** None

### WINDOW.ISFULLSCREEN
**Arguments:** None

### WINDOW.TOGGLEFULLSCREEN
**Arguments:** None

### WINDOW.ISRESIZED
**Arguments:** None

### WINDOW.SETTITLE
**Arguments:** `title$`
- `title$` (string): New window title

## Rendering Commands

### RENDER.CLEAR
**Arguments:** (multiple formats accepted)
- No arguments: Clears to black (0,0,0,255)
- `colorHandle`: Clear to color from heap handle
- `r, g, b`: Clear to RGB color (0-255)
- `r, g, b, a`: Clear to RGBA color (0-255)

### RENDER.FRAME
**Arguments:** None

### RENDER.SETBLEND / RENDER.SETBLENDMODE
**Arguments:** `mode`
- `mode` (int): Blend mode (use BLEND_* constants)

### RENDER.SETDEPTHWRITE / RENDER.SETDEPTHMASK
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable depth writing

### RENDER.SETDEPTHTEST
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable depth testing

### RENDER.SETSCISSOR
**Arguments:** `x, y, w, h`
- `x` (numeric): Scissor rectangle X position
- `y` (numeric): Scissor rectangle Y position
- `w` (numeric): Scissor rectangle width
- `h` (numeric): Scissor rectangle height

### RENDER.CLEARSCISSOR
**Arguments:** None

### RENDER.SETWIREFRAME
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable wireframe mode

### RENDER.SCREENSHOT
**Arguments:** `filename$`
- `filename$` (string): Screenshot filename

### RENDER.SETMSAA
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable MSAA

### RENDER.SETSHADOWMAPSIZE
**Arguments:** `size`
- `size` (int): Shadow map size in pixels

### RENDER.SETAMBIENT
**Arguments:** `r#, g#, b#` or `r#, g#, b#, a#`
- `r#, g#, b#` (float): RGB ambient color (0.0-1.0)
- `a#` (float, optional): Alpha multiplier (0.0-1.0)

### RENDER.SETMODE
**Arguments:** `mode`
- `mode` (int): Render mode

## World Management Commands

### WORLD.SETCENTER
**Arguments:** `x#, y#, z#`
- `x#, y#, z#` (float): World center coordinates

### WORLD.UPDATE
**Arguments:** `deltaTime#`
- `deltaTime#` (float): Time delta for update

### WORLD.STREAMENABLE
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable world streaming

### WORLD.PRELOAD
**Arguments:** `radius#`
- `radius#` (float): Preload radius

### WORLD.STATUS
**Arguments:** None

### WORLD.ISREADY
**Arguments:** None

### FOGMODE
**Arguments:** `mode%`
- `mode%` (int): Fog mode percentage

### FOGCOLOR
**Arguments:** `r, g, b`
- `r, g, b` (int): RGB fog color (0-255)

### FOGDENSITY
**Arguments:** `density#`
- `density#` (float): Fog density value

## Event/Automation Commands

### EVENT.LISTMAKE
**Arguments:** `path$`
- `path$` (string): Path for event list (can be empty)

### EVENT.LISTLOAD
**Arguments:** `path$`
- `path$` (string): Path to load event list from

### EVENT.LISTEXPORT
**Arguments:** `listHandle, path$`
- `listHandle` (handle): Event list handle
- `path$` (string): Export path

### EVENT.SETACTIVELIST
**Arguments:** `listHandle`
- `listHandle` (handle): Event list handle to set as active

### EVENT.RECSTART
**Arguments:** None

### EVENT.RECSTOP
**Arguments:** None

### EVENT.REPLAY
**Arguments:** `listHandle`
- `listHandle` (handle): Event list handle to replay

### EVENT.RECPLAYING
**Arguments:** None

### EVENT.ISPLAYING
**Arguments:** None

### EVENT.LISTCLEAR
**Arguments:** `listHandle`
- `listHandle` (handle): Event list handle to clear

### EVENT.LISTCOUNT
**Arguments:** `listHandle`
- `listHandle` (handle): Event list handle

### EVENT.LISTFREE
**Arguments:** `listHandle`
- `listHandle` (handle): Event list handle to free

## Decal Commands

### DECAL.MAKE
**Arguments:** `textureHandle, x#, y#, z#, width#, height#, lifetime#`
- `textureHandle` (handle): Texture handle for decal
- `x#, y#, z#` (float): 3D position
- `width#, height#` (float): Decal dimensions
- `lifetime#` (float): Decal lifetime in seconds

### DECAL.FREE
**Arguments:** `decalHandle`
- `decalHandle` (handle): Decal handle to free

### DECAL.SETPOS
**Arguments:** `decalHandle, x#, y#, z#`
- `decalHandle` (handle): Decal handle
- `x#, y#, z#` (float): New 3D position

### DECAL.SETSIZE
**Arguments:** `decalHandle, width#, height#`
- `decalHandle` (handle): Decal handle
- `width#, height#` (float): New dimensions

### DECAL.SETLIFETIME
**Arguments:** `decalHandle, lifetime#`
- `decalHandle` (handle): Decal handle
- `lifetime#` (float): New lifetime in seconds

### DECAL.DRAW
**Arguments:** `decalHandle`
- `decalHandle` (handle): Decal handle to draw

## Post-Processing Effect Commands

### EFFECT.SSAO
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable SSAO

### EFFECT.SSR
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable screen space reflections

### EFFECT.MOTIONBLUR
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable motion blur

### EFFECT.DEPTHOFFIELD
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable depth of field

### EFFECT.BLOOM
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable bloom effect

### EFFECT.TONEMAPPING
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable tone mapping

### EFFECT.SHARPEN
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable sharpen effect

### EFFECT.GRAIN
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable film grain

### EFFECT.VIGNETTE
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable vignette effect

### EFFECT.CHROMATICABERRATION
**Arguments:** `enabled`
- `enabled` (bool): Enable/disable chromatic aberration

## Render Target Commands

### RENDERTARGET.MAKE
**Arguments:** `width, height`
- `width` (int): Render target width
- `height` (int): Render target height

### RENDERTARGET.FREE
**Arguments:** `renderTargetHandle`
- `renderTargetHandle` (handle): Render target handle to free

### RENDERTARGET.BEGIN
**Arguments:** `renderTargetHandle`
- `renderTargetHandle` (handle): Render target handle to begin drawing to

### RENDERTARGET.END
**Arguments:** None

### RENDERTARGET.TEXTURE
**Arguments:** `renderTargetHandle`
- `renderTargetHandle` (handle): Render target handle

## Texture Commands

### TEXTURE.WIDTH
**Arguments:** `textureHandle`
- `textureHandle` (handle): Texture handle

### TEXTURE.HEIGHT
**Arguments:** `textureHandle`
- `textureHandle` (handle): Texture handle

### TEXTURE.SETFILTER
**Arguments:** `textureHandle, filterMode`
- `textureHandle` (handle): Texture handle
- `filterMode` (int): Filter mode (use TEXTURE_FILTER_* constants)

### TEXTURE.SETWRAP
**Arguments:** `textureHandle, wrapMode`
- `textureHandle` (handle): Texture handle
- `wrapMode` (int): Wrap mode (use TEXTURE_WRAP_* constants)

### TEXTURE.UPDATE
**Arguments:** `textureHandle, imageHandle`
- `textureHandle` (handle): Texture handle to update
- `imageHandle` (handle): Image handle with new data

## Global Shorthands (Easy Mode)

### SKYCOLOR
**Arguments:** Same as RENDER.CLEAR

### FPS
**Arguments:** None

### AMBIENTLIGHT
**Arguments:** None (currently no-op)

### SCREENWIDTH
**Arguments:** None

### SCREENHEIGHT
**Arguments:** None

---

## Argument Type Conventions

- **numeric**: Can be integer or float
- **int**: Must be integer value
- **float**: Must be floating-point value
- **string**: Must be string value
- **bool**: Must be boolean value (0/1)
- **handle**: Must be heap object handle
- **# suffix**: Indicates float value
- **% suffix**: Indicates integer percentage
- **$ suffix**: Indicates string value

## Constants

### Blend Modes
- `BLEND_ALPHA`
- `BLEND_ADDITIVE`
- `BLEND_MULTIPLIED`
- `BLEND_SUBTRACT`

### Window Flags
- `FLAG_VSYNC_HINT`
- `FLAG_FULLSCREEN_MODE`
- `FLAG_WINDOW_RESIZABLE`
- `FLAG_WINDOW_UNDECORATED`
- `FLAG_WINDOW_TRANSPARENT`
- `FLAG_WINDOW_HIDDEN`
- `FLAG_WINDOW_MINIMIZED`
- `FLAG_WINDOW_MAXIMIZED`
- `FLAG_WINDOW_UNFOCUSED`
- `FLAG_WINDOW_TOPMOST`
- `FLAG_WINDOW_ALWAYS_RUN`
- `FLAG_WINDOW_HIGHDPI`
- `FLAG_WINDOW_INTERLACED`

### Texture Filter Modes
- `TEXTURE_FILTER_POINT`
- `TEXTURE_FILTER_BILINEAR`
- `TEXTURE_FILTER_TRILINEAR`
- `TEXTURE_FILTER_ANISOTROPIC_4X`
- `TEXTURE_FILTER_ANISOTROPIC_8X`
- `TEXTURE_FILTER_ANISOTROPIC_16X`

### Texture Wrap Modes
- `TEXTURE_WRAP_REPEAT`
- `TEXTURE_WRAP_CLAMP`
- `TEXTURE_WRAP_MIRROR_REPEAT`
- `TEXTURE_WRAP_MIRROR_CLAMP`

---

## Usage Examples

### Basic Window Setup
```basic
' Open window
WINDOW.OPEN(800, 600, "My Game")

' Main game loop
WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(50, 50, 50, 255)
    
    ' Your game logic here
    
    RENDER.FRAME()
WEND

' Clean up
WINDOW.CLOSE()
```

### Audio Playback
```basic
' Load and play music
music = AUDIO.LOADMUSIC("background.mp3")
AUDIO.UPDATEMUSIC(music)

' Load and play sound
sound = AUDIO.LOADSOUND("explosion.wav")
' Play sound would be implemented separately
```

### Window Configuration
```basic
' Set window properties
WINDOW.SETFPS(60)
WINDOW.SETTITLE("My Application")
WINDOW.SETPOSITION(100, 100)
WINDOW.SETSIZE(1024, 768)

' Check window state
width = WINDOW.WIDTH()
height = WINDOW.HEIGHT()
fps = WINDOW.GETFPS()
```

### Rendering Setup
```basic
' Configure rendering
RENDER.SETBLEND(BLEND_ALPHA)
RENDER.SETAMBIENT(0.2, 0.2, 0.2)
RENDER.SETDEPTHWRITE(1)
RENDER.SETDEPTHTEST(1)

' Clear screen with custom color
RENDER.CLEAR(128, 128, 255, 255)
```

This documentation covers all currently implemented MoonBasic commands with their exact argument specifications.
