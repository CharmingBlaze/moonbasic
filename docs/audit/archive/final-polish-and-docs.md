Systems for the programmer to build
1. APP system

Handles the application, window, timing, and loop state.

APP.OPEN(1280, 720, "MoonBasic Game")
APP.CLOSE()
APP.SHOULDCLOSE()
APP.SETFPS(60)
APP.GETFPS()
APP.WIDTH()
APP.HEIGHT()
APP.TIME()
APP.DELTA()
APP.VERSION()

Needed internally:

window creation
frame timing
fixed update support
fullscreen/windowed mode
window resize
exit state
FPS cap
delta time
2. RENDER system

Handles frame drawing and render modes.

RENDER.CLEAR(20, 20, 30)
RENDER.BEGIN()
RENDER.END()
RENDER.FRAME()
RENDER.SETBACKGROUND(20, 20, 30)
RENDER.SETWIREFRAME(true)
RENDER.SCREENSHOT("shot.png")

Needed internally:

main render loop
clear color
2D pass
3D pass
debug pass
UI pass
wireframe mode
screenshot capture
render ordering
3. SCENE system

Handles world containers, loading, saving, and active scene.

scene = SCENE.CREATE("Level1")
SCENE.SETACTIVE(scene)
SCENE.SAVE(scene, "levels/level1.scene")
SCENE.LOAD("levels/level1.scene")
SCENE.CLEAR()
SCENE.FIND("Player")
SCENE.COUNT()

Needed internally:

scene registry
active scene
entity list
scene serialization
scene loading
scene reset
persistent entities
scene metadata
4. ENTITY system

This is the main object system. Everything should be an entity.

player = ENTITY.CREATE("Player")
cube = ENTITY.CREATECUBE("Crate", 2)
sphere = ENTITY.CREATESPHERE("Ball", 1)
pivot = ENTITY.CREATEPIVOT("Pivot")
ENTITY.DESTROY(player)
ENTITY.CLONE(player)

Transform commands:

ENTITY.SETPOSITION(player, 0, 1, 5)
ENTITY.SETROTATION(player, 0, 90, 0)
ENTITY.SETSCALE(player, 1, 1, 1)

ENTITY.MOVE(player, 0, 0, 1)
ENTITY.TURN(player, 0, 1, 0)

x = ENTITY.GETX(player)
y = ENTITY.GETY(player)
z = ENTITY.GETZ(player)

ENTITY.LOOKAT(player, target)
ENTITY.POINTAT(player, 0, 0, 0)

Hierarchy commands:

ENTITY.SETPARENT(child, parent)
ENTITY.CLEARPARENT(child)
ENTITY.GETPARENT(child)
ENTITY.CHILDCOUNT(parent)
ENTITY.GETCHILD(parent, index)

State commands:

ENTITY.SHOW(player)
ENTITY.HIDE(player)
ENTITY.SETVISIBLE(player, true)
ENTITY.ISVISIBLE(player)
ENTITY.SETNAME(player, "Player")
ENTITY.GETNAME(player)
ENTITY.SETTAG(player, "enemy")
ENTITY.HASTAG(player, "enemy")

Needed internally:

entity ids
entity registry
transform component
parent/child hierarchy
local transform
world transform
dirty transform update
visibility
names
tags
layers
clone/destroy safety
5. CAMERA system

Cameras are entities, but have camera-specific commands.

camera = CAMERA.CREATE("MainCamera")
CAMERA.SETACTIVE(camera)
CAMERA.SETFOV(camera, 70)
CAMERA.SETNEAR(camera, 0.1)
CAMERA.SETFAR(camera, 1000)
CAMERA.LOOKAT(camera, 0, 0, 0)
CAMERA.FOLLOW(camera, player, 0, 3, -8)

View/projection commands:

CAMERA.BEGIN(camera)
CAMERA.END(camera)
CAMERA.SCREENRAY(camera, mouseX, mouseY)
CAMERA.WORLDTOSCREEN(camera, x, y, z)
CAMERA.SCREENTOWORLD(camera, x, y)

Needed internally:

active camera
perspective camera
orthographic camera
camera entity transform
FOV
near/far planes
follow camera helper
screen ray generation
world/screen conversion
6. LIGHT system

Lights are entities too.

light = LIGHT.CREATEPOINT("Lamp")
sun = LIGHT.CREATEDIRECTIONAL("Sun")
spot = LIGHT.CREATESPOT("Torch")

LIGHT.SETCOLOR(light, 255, 220, 180)
LIGHT.SETINTENSITY(light, 2.0)
LIGHT.SETRANGE(light, 15)
LIGHT.SETDIRECTION(sun, -1, -2, -1)
LIGHT.SETSHADOWS(sun, true)

Needed internally:

directional lights
point lights
spot lights
ambient light
light uniforms
shadow flag placeholder
lighting shader support
7. MESH system

Handles mesh creation and mesh data.

mesh = MESH.CREATE()
cubeMesh = MESH.CUBE(2, 2, 2)
sphereMesh = MESH.SPHERE(1, 16, 16)
planeMesh = MESH.PLANE(10, 10)
cylinderMesh = MESH.CYLINDER(1, 2, 16)

ENTITY.SETMESH(player, cubeMesh)

Editable mesh commands:

MESH.ADDVERTEX(mesh, x, y, z)
MESH.ADDFACE(mesh, v1, v2, v3)
MESH.ADDQUAD(mesh, v1, v2, v3, v4)
MESH.RECALCNORMALS(mesh)
MESH.UPLOAD(mesh)

Needed internally:

mesh registry
primitive generation
vertex buffer
index buffer
normals
uvs
colors
upload/update to GPU
editable mesh support
8. MODEL system

Handles loaded models.

model = MODEL.LOAD("assets/player.glb")
ENTITY.SETMODEL(player, model)
MODEL.UNLOAD(model)
MODEL.ANIMCOUNT(model)
MODEL.GETANIMNAME(model, 0)

Needed internally:

GLB/GLTF loading
OBJ loading
model registry
model materials
model animation list
model unload
missing asset fallback
9. MATERIAL system

Handles material appearance.

mat = MATERIAL.CREATE("PlayerMat")
MATERIAL.SETCOLOR(mat, 255, 255, 255, 255)
MATERIAL.SETTEXTURE(mat, texture)
MATERIAL.SETALPHA(mat, 0.5)
MATERIAL.SETMETALLIC(mat, 0.0)
MATERIAL.SETROUGHNESS(mat, 0.8)

ENTITY.SETMATERIAL(player, mat)

Needed internally:

material registry
base color
texture slot
alpha
unlit mode
lit mode
toon mode
PBR-lite values
shader assignment
10. TEXTURE system

Handles image loading and texture management.

tex = TEXTURE.LOAD("assets/crate.png")
TEXTURE.UNLOAD(tex)
TEXTURE.WIDTH(tex)
TEXTURE.HEIGHT(tex)

ENTITY.SETTEXTURE(player, tex)

Needed internally:

texture registry
image loading
missing texture fallback
transparent textures
texture filtering
texture wrapping
texture unload
11. ANIMATION system

Handles animated models and entity animation state.

ANIM.PLAY(player, "Run")
ANIM.STOP(player)
ANIM.PAUSE(player)
ANIM.RESUME(player)
ANIM.SETTIME(player, 0.5)
ANIM.SETSPEED(player, 1.0)
ANIM.ISPLAYING(player)
ANIM.CURRENT(player)

Needed internally:

animation clips
current animation
animation speed
loop mode
play/pause/stop
skeletal animation update
animation events later
12. INPUT system

Handles keyboard, mouse, and gamepads.

INPUT.KEYDOWN(KEY_W)
INPUT.KEYHIT(KEY_SPACE)
INPUT.KEYUP(KEY_ESCAPE)

INPUT.MOUSEDOWN(MOUSE_LEFT)
INPUT.MOUSEHIT(MOUSE_LEFT)
INPUT.MOUSEX()
INPUT.MOUSEY()
INPUT.MOUSEDELTA_X()
INPUT.MOUSEDELTA_Y()
INPUT.MOUSEWHEEL()

INPUT.GAMEPADCONNECTED(0)
INPUT.GAMEPADBUTTONDOWN(0, PAD_A)
INPUT.GAMEPADAXIS(0, PAD_LEFT_X)

Needed internally:

keyboard state
mouse state
gamepad state
pressed/released detection
input aliases
input mapping later
13. ACTION system

This sits above raw input. Users map controls to actions.

ACTION.BINDKEY("Jump", KEY_SPACE)
ACTION.BINDKEY("Forward", KEY_W)
ACTION.BINDGAMEPAD("Jump", PAD_A)

ACTION.DOWN("Forward")
ACTION.HIT("Jump")
ACTION.VALUE("MoveX")

Needed internally:

action map
keyboard bindings
mouse bindings
gamepad bindings
axis bindings
save/load bindings
14. PHYSICS system

General physics world control.

PHYSICS.CREATEWORLD()
PHYSICS.SETGRAVITY(0, -9.8, 0)
PHYSICS.STEP()
PHYSICS.DEBUGDRAW(true)

Needed internally:

physics world
fixed timestep
gravity
debug drawing
body registry
sync entity to body
sync body to entity
15. BODY system

Physics bodies attached to entities.

BODY.ADDSTATICBOX(wall, 2, 2, 2)
BODY.ADDDYNAMICBOX(crate, 1, 1, 1)
BODY.ADDSPHERE(ball, 1)
BODY.ADDCAPSULE(player, 0.4, 1.8)

BODY.SETMASS(crate, 10)
BODY.SETFRICTION(crate, 0.5)
BODY.SETBOUNCE(crate, 0.2)
BODY.APPLYFORCE(crate, 0, 10, 0)
BODY.APPLYIMPULSE(crate, 0, 5, 0)

Needed internally:

static bodies
dynamic bodies
kinematic bodies
box shape
sphere shape
capsule shape
mesh collider later
mass
friction
restitution
force/impulse
entity transform sync
16. COLLISION system

Beginner-friendly collision separate from full physics.

COLLISION.SETRADIUS(player, 0.4)
COLLISION.SETBOX(wall, 2, 2, 2)
COLLISION.SETTYPE(player, 1)
COLLISION.SETTYPE(wall, 2)
COLLISION.RULE(1, 2, COLLIDE_SLIDE)

COLLISION.UPDATE()
COLLISION.HIT(player, wall)
COLLISION.COUNT(player)

Needed internally:

simple collision shapes
collision groups
collision rules
overlap tests
sliding collision
collision events
last collision data
17. PICK system

Ray picking and mouse picking.

hit = PICK.MOUSE(camera)
hit = PICK.RAY(x, y, z, dx, dy, dz)
hit = PICK.ENTITY(player, camera)

PICK.HIT()
PICK.ENTITY()
PICK.X()
PICK.Y()
PICK.Z()
PICK.NX()
PICK.NY()
PICK.NZ()
PICK.DISTANCE()

Needed internally:

ray generation
ray vs box
ray vs sphere
ray vs mesh
nearest hit
hit point
hit normal
hit entity id
layer mask
18. AUDIO system

Handles sound effects and music.

sound = AUDIO.LOADSOUND("jump.wav")
AUDIO.PLAYSOUND(sound)
AUDIO.STOPSOUND(sound)
AUDIO.SETVOLUME(sound, 0.8)

music = AUDIO.LOADMUSIC("theme.ogg")
AUDIO.PLAYMUSIC(music)
AUDIO.UPDATEMUSIC(music)
AUDIO.STOPMUSIC(music)

Needed internally:

sound registry
music registry
one-shot sound
looping sound
streaming music
volume
pitch
pan
audio cleanup
19. AUDIO3D system

Sound attached to world positions/entities.

snd = AUDIO3D.LOAD("explosion.wav")
AUDIO3D.PLAYAT(snd, 10, 0, 5)
AUDIO3D.ATTACH(snd, enemy)
AUDIO3D.SETLISTENER(camera)
AUDIO3D.SETRANGE(snd, 20)

Needed internally:

listener entity
3D sound positions
distance falloff
attached sound emitters
doppler later
20. UI system

Simple immediate UI for tools, debug screens, and game menus.

UI.BEGIN()
UI.LABEL("Health: " + health)
if UI.BUTTON("Start") then
    StartGame()
endif
UI.SLIDER("Volume", volume, 0, 1)
UI.END()

Needed internally:

buttons
labels
sliders
checkboxes
text input
panels
debug overlay
font rendering
mouse interaction
21. FONT/TEXT system

Text drawing separate from UI.

font = FONT.LOAD("assets/font.ttf", 32)
TEXT.DRAW("Hello", 20, 20)
TEXT.DRAWFONT(font, "Score: 100", 20, 20)
TEXT.SIZE("Hello")

Needed internally:

default font
TTF loading
text measurement
font registry
2D text drawing
optional 3D text later
22. SPRITE system

2D game and UI sprite support.

sprite = SPRITE.CREATE(texture)
SPRITE.SETPOSITION(sprite, 100, 200)
SPRITE.SETROTATION(sprite, 45)
SPRITE.SETSCALE(sprite, 2, 2)
SPRITE.DRAW(sprite)

Needed internally:

sprite registry
texture assignment
source rectangle
origin/pivot
rotation
scale
color tint
layer/depth sorting
23. TILEMAP system

For 2D games.

map = TILEMAP.LOAD("level.tmx")
TILEMAP.DRAW(map)
TILEMAP.GETTILE(map, layer, x, y)
TILEMAP.SETTILE(map, layer, x, y, tileId)

Needed internally:

tilemap loading
tile layers
collision layers
tileset textures
camera culling
TMX or simple JSON support
24. TERRAIN system

For simple 3D worlds.

terrain = TERRAIN.CREATE(128, 128)
TERRAIN.LOADHEIGHTMAP(terrain, "height.png")
TERRAIN.SETTEXTURE(terrain, texture)
TERRAIN.SETHEIGHT(terrain, x, z, h)
TERRAIN.GETHEIGHT(terrain, x, z)

Needed internally:

heightmap mesh
terrain collision
texture layers
height queries
terrain chunks later
25. PARTICLE system

For fire, smoke, magic, sparks.

fx = PARTICLE.CREATE()
PARTICLE.SETTEXTURE(fx, texture)
PARTICLE.SETRATE(fx, 50)
PARTICLE.SETLIFETIME(fx, 2.0)
PARTICLE.SETSPEED(fx, 1, 5)
PARTICLE.PLAY(fx)
PARTICLE.STOP(fx)

Needed internally:

particle emitter
spawn rate
lifetime
velocity
gravity
color over time
size over time
billboarding
26. SCRIPT/TIMER system

Useful beginner helpers.

TIMER.AFTER(1.0, "SpawnEnemy")
TIMER.EVERY(0.5, "Shoot")
TIMER.CANCEL(timerId)

Needed internally:

delayed callbacks
repeating callbacks
timer registry
frame update
27. SAVE system

Save/load game state.

SAVE.SET("level", 3)
SAVE.SET("health", 80)
SAVE.WRITE("save1.json")

SAVE.READ("save1.json")
level = SAVE.GET("level")

Needed internally:

key/value save data
JSON save file
basic types
arrays later
safe file paths
28. ASSET system

This is very important.

ASSET.LOADPACK("assets/assets.json")

playerTex = ASSET.TEXTURE("player")
playerModel = ASSET.MODEL("player")
jumpSound = ASSET.SOUND("jump")

Needed internally:

asset manifest
asset id lookup
texture/model/audio loading
missing asset fallback
asset unloading
hot reload later

Example assets.json:

{
  "textures": {
    "player": "textures/player.png",
    "crate": "textures/crate.png"
  },
  "models": {
    "hero": "models/hero.glb"
  },
  "sounds": {
    "jump": "audio/jump.wav"
  }
}
29. RESOURCE system

Low-level file helpers.

FILE.EXISTS("data/config.json")
FILE.READTEXT("data/config.json")
FILE.WRITETEXT("data/save.json", text)
FILE.DELETE("data/save.json")

Needed internally:

safe paths
read text
write text
read bytes
write bytes
directory create
file exists
30. JSON system

Needed for saves, configs, asset files.

json = JSON.PARSE(text)
JSON.GET(json, "player.health")
JSON.SET(json, "player.health", 100)
text = JSON.STRINGIFY(json)

Needed internally:

JSON parse
JSON stringify
path lookup
numbers
strings
booleans
arrays
objects
31. MATH system

Game math helpers.

MATH.RAND(1, 10)
MATH.RANDF(0, 1)
MATH.CLAMP(value, 0, 100)
MATH.LERP(a, b, t)
MATH.DISTANCE(x1, y1, z1, x2, y2, z2)

Needed internally:

random numbers
vectors
quaternions
matrices
lerp
clamp
distance
angle helpers
32. VEC3 system

Cleaner vector operations.

v = VEC3.CREATE(0, 1, 5)
v = VEC3.ADD(a, b)
v = VEC3.NORMALIZE(v)
len = VEC3.LENGTH(v)

Needed internally:

vector type
add/subtract
multiply
normalize
dot
cross
length
distance
33. DEBUG system

This is needed before VS Code.

DEBUG.LOG("Player spawned")
DEBUG.WARN("Missing texture")
DEBUG.ERROR("Bad state")
DEBUG.DRAWLINE(0, 0, 0, 10, 0, 0)
DEBUG.DRAWBOX(player)
DEBUG.SHOWFPS(true)
DEBUG.SHOWSTATS(true)

Needed internally:

logging
warnings
runtime error display
FPS counter
draw line
draw box
draw sphere
draw ray
memory/object stats
34. ERROR system

Better compiler/runtime errors.

Needed features:

file name
line number
column number
function name
call stack
bad command name
bad argument count
bad type explanation
suggest closest command
runtime crash recovery

Example:

main.mb:42:5
ENTITY.SETPOSITON(player, 0, 1, 5)
      ^^^^^^^^^^

Unknown command: ENTITY.SETPOSITON

Did you mean:
  ENTITY.SETPOSITION
35. PROJECT system

Create and run projects.

moonbasic new MyGame
moonbasic run
moonbasic build
moonbasic package

Needed internally:

project manifest
main file
asset folder
build folder
output folder
template creation
run current project
release package

Example:

{
  "name": "MyGame",
  "main": "main.mb",
  "assets": "assets",
  "version": "0.1.0"
}
36. PACKAGE/EXPORT system

This is how users share games.

moonbasic package windows
moonbasic package linux
moonbasic package web

Needed internally:

copy runtime
compile bytecode
copy assets
write config
zip release
include license/readme

Start with Windows only.

37. MODULE system

For splitting code.

IMPORT "player.mb"
IMPORT "enemy.mb"
IMPORT "ui/hud.mb"

Needed internally:

import resolver
relative paths
module cache
duplicate import protection
public/private later
38. COMMAND HELP system

Built-in help before VS Code.

HELP("ENTITY.SETPOSITION")
HELP("CAMERA")

Needed internally:

command registry
command descriptions
argument list
examples
search help
print to console
39. TEST system

For testing language and runtime.

moonbasic test

Needed internally:

unit tests for compiler
runtime smoke tests
golden output tests
example compile tests
example run tests
binding tests
40. TEMPLATE system

Starter projects.

moonbasic new My3DGame --template 3d
moonbasic new MyPlatformer --template platformer
moonbasic new MyMenu --template ui

Templates needed:

empty
2d
3d
first-person
third-person
platformer
top-down
physics
ui-menu
Build order

Tell the programmer to build in this order:

1. APP
2. RENDER
3. SCENE
4. ENTITY
5. CAMERA
6. LIGHT
7. MESH
8. MATERIAL
9. TEXTURE
10. INPUT
11. DEBUG
12. ERROR
13. ASSET
14. MODEL
15. ANIMATION
16. PICK
17. COLLISION
18. PHYSICS
19. BODY
20. AUDIO
21. SPRITE
22. UI
23. PROJECT
24. PACKAGE
25. COMMAND HELP
26. TEMPLATES

Do not start with VS Code.

Do not start with advanced networking.

Do not start with a visual editor.

Get this working first:

APP.OPEN(1280, 720, "Test")

camera = CAMERA.CREATE("Main Camera")
CAMERA.SETACTIVE(camera)
ENTITY.SETPOSITION(camera, 0, 2, -8)
CAMERA.LOOKAT(camera, 0, 0, 0)

cube = ENTITY.CREATECUBE("Cube", 2)
ENTITY.SETPOSITION(cube, 0, 0, 5)

light = LIGHT.CREATEPOINT("Light")
ENTITY.SETPOSITION(light, 0, 5, -3)

WHILE NOT APP.SHOULDCLOSE()
    ENTITY.TURN(cube, 0, 60 * APP.DELTA(), 0)

    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN()
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

That is the foundation.