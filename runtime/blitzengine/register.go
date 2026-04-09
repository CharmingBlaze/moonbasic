package blitzengine

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const (
	raylibFlagVsyncHint  = 0x00000040
	raylibFlagMsaa4xHint = 0x00000020
)

// registerBlitzAPI wires Blitz-style flat names. See docs/reference/moonbasic-command-set/.
func registerBlitzAPI(m *Module, reg runtime.Registrar) {
	reg.Register("APPTITLE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"WINDOW.SETTITLE", args...)
	})
	reg.Register("SETFPS", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"WINDOW.SETFPS", args...)
	})
	reg.Register("DELTATIME", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("DELTATIME expects 0 arguments")
		}
		return call(rt,"TIME.DELTA")
	})
	reg.Register("TIMEMS", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("TIMEMS expects 0 arguments")
		}
		return call(rt,"TICKCOUNT")
	})
	// SLEEP is registered by core — do not override (would recurse).
	reg.Register("FINISH", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("FINISH expects 0 arguments (ends program; END is a reserved keyword)")
		}
		return call(rt,"ENDGAME")
	})

	reg.Register("GRAPHICS", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("GRAPHICS expects 3 arguments (width, height, depth) — depth is informational; title is empty")
		}
		title := value.FromStringIndex(rt.Heap.Intern(""))
		return call(rt,"WINDOW.OPEN", args[0], args[1], title)
	})
	reg.Register("GRAPHICS3D", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("GRAPHICS3D expects 3 arguments (width, height, depth)")
		}
		title := value.FromStringIndex(rt.Heap.Intern(""))
		if _, err := call(rt,"WINDOW.OPEN", args[0], args[1], title); err != nil {
			return value.Nil, err
		}
		dep, _ := args[2].ToInt()
		if dep >= 24 {
			if _, err := call(rt,"WINDOW.SETFLAG", value.FromInt(raylibFlagMsaa4xHint)); err != nil {
				return value.Nil, err
			}
		}
		return value.Nil, nil
	})
	reg.Register("SETVSYNC", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("SETVSYNC expects 1 argument (non-zero = enable vsync hint)")
		}
		on, _ := args[0].ToInt()
		if on != 0 {
			return call(rt,"WINDOW.SETFLAG", value.FromInt(raylibFlagVsyncHint))
		}
		return call(rt,"WINDOW.CLEARFLAG", value.FromInt(raylibFlagVsyncHint))
	})
	reg.Register("SETCLEARCOLOR", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("SETCLEARCOLOR expects 3 arguments (r, g, b) 0–255")
		}
		return call(rt,"RENDER.CLEAR", args[0], args[1], args[2], value.FromInt(255))
	})
	reg.Register("CLEAR", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CLEAR expects 0 arguments")
		}
		return call(rt,"RENDER.CLEAR")
	})
	reg.Register("FLIP", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("FLIP expects 0 arguments")
		}
		return call(rt,"RENDER.FRAME")
	})

	reg.Register("SETCOLOR", "blitzengine", m.setColor)
	reg.Register("SETALPHA", "blitzengine", m.setAlpha)
	reg.Register("SETORIGIN", "blitzengine", m.setOrigin)
	reg.Register("SETVIEWPORT", "blitzengine", m.setViewport)
	reg.Register("PLOT", "blitzengine", m.plot)
	reg.Register("LINE", "blitzengine", m.line)
	reg.Register("RECT", "blitzengine", m.rect)
	reg.Register("OVAL", "blitzengine", m.oval)
	reg.Register("TEXT", "blitzengine", m.textDraw)

	reg.Register("CREATEWORLD", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CREATEWORLD expects 0 arguments (implicit after WINDOW.OPEN)")
		}
		return value.Nil, nil
	})
	reg.Register("UPDATEWORLD", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.UPDATE", args...)
	})
	// UpdatePhysics — one call per frame: ENTITY.UPDATE(dt), optional WORLD/2D/3D physics (errors ignored if inactive).
	reg.Register("UPDATEPHYSICS", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("UPDATEPHYSICS expects 0 arguments")
		}
		dt, err := call(rt, "TIME.DELTA")
		if err != nil {
			return value.Nil, err
		}
		if _, err := call(rt, "ENTITY.UPDATE", dt); err != nil {
			return value.Nil, err
		}
		_, _ = call(rt, "WORLD.UPDATE", dt)
		_, _ = call(rt, "PHYSICS2D.STEP")
		_, _ = call(rt, "PHYSICS3D.STEP", dt)
		return value.Nil, nil
	})
	reg.Register("RENDERWORLD", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("RENDERWORLD expects 0 arguments")
		}
		return call(rt,"ENTITY.DRAWALL")
	})
	reg.Register("CLEARWORLD", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("CLEARWORLD expects 0 arguments")
		}
		return call(rt,"ENTITY.CLEARSCENE")
	})
	reg.Register("SETAMBIENT", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"RENDER.SETAMBIENT", args...)
	})
	reg.Register("SETFOG", "blitzengine", m.setFog)
	reg.Register("SETWIREFRAME", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("SETWIREFRAME expects 1 argument (non-zero = wireframe)")
		}
		on, _ := args[0].ToInt()
		return call(rt,"RENDER.SETWIREFRAME", value.FromBool(on != 0))
	})

	reg.Register("CREATECAMERA", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) > 1 {
			return value.Nil, fmt.Errorf("CREATECAMERA expects 0 or 1 argument (parent — optional, ignored)")
		}
		return call(rt,"CAMERA.MAKE")
	})

	reg.Register("MOVESTEPX", "blitzengine", runtime.AdaptLegacy(m.entMoveStepX))
	reg.Register("MOVESTEPZ", "blitzengine", runtime.AdaptLegacy(m.entMoveStepZ))
	reg.Register("DIST3D", "blitzengine", runtime.AdaptLegacy(m.dist3D))
	reg.Register("COLORPRINT", "blitzengine", m.colorPrint)
	reg.Register("FPS", "blitzengine", m.fps)
	reg.Register("MILLISECS", "blitzengine", m.milliSecs)
	reg.Register("SCREENWIDTH", "blitzengine", m.screenWidth)
	reg.Register("SCREENHEIGHT", "blitzengine", m.screenHeight)
	reg.Register("POSITIONCAMERA", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"CAMERA.SETPOS", args...)
	})
	reg.Register("ROTATECAMERA", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"CAMERA.ROTATE", args...)
	})
	reg.Register("MOVECAMERA", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"CAMERA.MOVE", args...)
	})
	reg.Register("CAMERARANGE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"CAMERA.SETRANGE", args...)
	})
	reg.Register("CAMERAZOOM", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"CAMERA.ZOOM", args...)
	})
	reg.Register("CAMERAVIEWPORT", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("CAMERAVIEWPORT expects 5 arguments (cam, x, y, w, h) — cam ignored; sets scissor")
		}
		return call(rt,"RENDER.SETSCISSOR", args[1], args[2], args[3], args[4])
	})
	reg.Register("CAMERAPROJECT", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"CAMERA.WORLDTOSCREEN", args...)
	})
	reg.Register("CAMERAPICK", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"CAMERA.PICK", args...)
	})

	reg.Register("CREATELIGHT", "blitzengine", m.createLight)
	reg.Register("LIGHTCOLOR", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"LIGHT.SETCOLOR", args...)
	})
	reg.Register("LIGHTRANGE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"LIGHT.SETRANGE", args...)
	})
	reg.Register("LIGHTCONE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("LIGHTCONE expects 3 arguments (light, inner#, outer#)")
		}
		if _, err := call(rt,"LIGHT.SETINNERCONE", args[0], args[1]); err != nil {
			return value.Nil, err
		}
		return call(rt,"LIGHT.SETOUTERCONE", args[0], args[2])
	})
	reg.Register("LIGHTPOSITION", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"LIGHT.SETPOSITION", args...)
	})
	reg.Register("LIGHTPOINTAT", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("LIGHTPOINTAT expects 4 arguments (light, x#, y#, z#)")
		}
		return call(rt,"LIGHT.SETTARGET", args[0], args[1], args[2], args[3])
	})

	reg.Register("CREATECUBE", "blitzengine", m.createCube)
	reg.Register("CREATESPHERE", "blitzengine", m.createSphereBlitz)
	reg.Register("CREATEPLANE", "blitzengine", m.createPlaneBlitz)
	reg.Register("CREATEMESH", "blitzengine", m.createMeshBlitz)
	reg.Register("CREATESPRITE3D", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("CREATESPRITE3D is not implemented — use ENTITY.CREATECUBE + ENTITY.TEXTURE")
	})
	reg.Register("COPYENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.COPY", args...)
	})
	reg.Register("FREEENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.FREE", args...)
	})
	reg.Register("POSITIONENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.POSITIONENTITY", args...)
	})
	reg.Register("ROTATEENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ROTATEENTITY", args...)
	})
	reg.Register("SCALEENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.SCALE", args...)
	})
	reg.Register("MOVEENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.MOVEENTITY", args...)
	})
	reg.Register("TURNENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.TURNENTITY", args...)
	})
	reg.Register("POINTENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.POINTENTITY", args...)
	})
	reg.Register("ALIGNTOVECTOR", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ALIGNTOVECTOR", args...)
	})
	reg.Register("ENTITYX", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ENTITYX", args...)
	})
	reg.Register("ENTITYY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ENTITYY", args...)
	})
	reg.Register("ENTITYZ", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ENTITYZ", args...)
	})
	reg.Register("ENTITYPITCH", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ENTITYPITCH", args...)
	})
	reg.Register("ENTITYYAW", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ENTITYYAW", args...)
	})
	reg.Register("ENTITYROLL", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ENTITYROLL", args...)
	})
	reg.Register("ENTITYSCALEX", "blitzengine", m.entityScaleX)
	reg.Register("ENTITYSCALEY", "blitzengine", m.entityScaleY)
	reg.Register("ENTITYSCALEZ", "blitzengine", m.entityScaleZ)
	reg.Register("SHOWENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.SHOW", args...)
	})
	reg.Register("HIDEENTITY", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.HIDE", args...)
	})
	reg.Register("ENTITYALPHA", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ALPHA", args...)
	})
	reg.Register("ENTITYCOLOR", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.COLOR", args...)
	})
	reg.Register("ENTITYFX", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.FX", args...)
	})
	reg.Register("ENTITYORDER", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.ORDER", args...)
	})
	reg.Register("ENTITYPARENT", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.PARENT", args...)
	})

	reg.Register("LOADMESH", "blitzengine", m.loadMeshParent)
	reg.Register("LOADANIMMESH", "blitzengine", m.loadAnimMeshParent)
	reg.Register("MESHWIDTH", "blitzengine", m.meshWidth)
	reg.Register("MESHHEIGHT", "blitzengine", m.meshHeight)
	reg.Register("MESHDEPTH", "blitzengine", m.meshDepth)
	reg.Register("ADDSURFACE", "blitzengine", notImplMesh)
	reg.Register("ADDVERTEX", "blitzengine", notImplMesh)
	reg.Register("ADDTRIANGLE", "blitzengine", notImplMesh)
	reg.Register("VERTEXX", "blitzengine", notImplMesh)
	reg.Register("VERTEXY", "blitzengine", notImplMesh)
	reg.Register("VERTEXZ", "blitzengine", notImplMesh)

	reg.Register("CREATETEXTURE", "blitzengine", m.createTexture)
	reg.Register("SCALETEXTURE", "blitzengine", notImplTexUV)
	reg.Register("ROTATETEXTURE", "blitzengine", notImplTexUV)
	reg.Register("TEXTURECOORDS", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"TEXTURE.SETWRAP", args...)
	})
	reg.Register("ENTITYTEXTURE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"ENTITY.TEXTURE", args...)
	})

	reg.Register("CREATESPRITE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"SPRITE.LOAD", args...)
	})
	reg.Register("SPRITE", "blitzengine", m.spriteAt)
	reg.Register("MOVESPRITE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"SPRITE.SETPOS", args...)
	})
	reg.Register("SPRITEIMAGE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("SPRITEIMAGE: use SPRITE.LOAD to change texture")
	})
	reg.Register("SPRITECOLOR", "blitzengine", m.spriteNoOpTint)
	reg.Register("SPRITEALPHA", "blitzengine", m.spriteNoOpTint)
	reg.Register("SPRITEHIT", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"SPRITE.HIT", args...)
	})

	reg.Register("LOADSOUND", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"AUDIO.LOADSOUND", args...)
	})
	reg.Register("PLAYSOUND", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"AUDIO.PLAY", args...)
	})
	reg.Register("LOOPSOUND", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"AUDIO.PLAY", args...)
	})
	reg.Register("STOPSOUND", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"AUDIO.STOP", args...)
	})
	reg.Register("SOUNDVOLUME", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"AUDIO.SETSOUNDVOLUME", args...)
	})
	reg.Register("SOUNDPAN", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"AUDIO.SETSOUNDPAN", args...)
	})
	reg.Register("SOUNDPITCH", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"AUDIO.SETSOUNDPITCH", args...)
	})

	reg.Register("MOUSEZ", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("MOUSEZ expects 0 arguments")
		}
		return call(rt,"INPUT.MOUSEWHEELMOVE")
	})
	reg.Register("MOUSEDOWN", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"INPUT.MOUSEDOWN", args...)
	})
	reg.Register("MOUSEHIT", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"INPUT.MOUSEHIT", args...)
	})
	reg.Register("MOVEMOUSE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"INPUT.SETMOUSEPOS", args...)
	})
	reg.Register("FLUSHKEYS", "blitzengine", flushNoOp)
	reg.Register("FLUSHMOUSE", "blitzengine", flushNoOp)

	reg.Register("READFILE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"FILE.READALLTEXT", args...)
	})
	// Spec: whole-file write. Replaces legacy flat WRITEFILE (stream) — use FILE.WRITE(handle, …) for streaming.
	reg.Register("WRITEFILE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"FILE.WRITEALLTEXT", args...)
	})
	reg.Register("CLOSEFILE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"FILE.CLOSE", args...)
	})
	reg.Register("READLINE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"FILE.READLINE", args...)
	})
	reg.Register("WRITELINE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"FILE.WRITELN", args...)
	})
	reg.Register("FILEEXISTS", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"FILE.EXISTS", args...)
	})
	reg.Register("DELETEFILE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"UTIL.DELETEFILE", args...)
	})
	reg.Register("COPYFILE", "blitzengine", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return call(rt,"UTIL.COPYFILE", args...)
	})
}

func flushNoOp(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("expects 0 arguments")
	}
	return value.Nil, nil
}

func notImplMesh(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("mesh surface API not implemented — use MESH.MAKECUSTOM / ENTITY.LOADMESH (see MESH.md)")
}

func notImplTexUV(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("not implemented — use TEXTURE.SETFILTER / DRAW.TEXTUREPRO for UV transforms")
}

func valFloat(v value.Value) float64 {
	if f, ok := v.ToFloat(); ok {
		return f
	}
	if i, ok := v.ToInt(); ok {
		return float64(i)
	}
	return 0
}

func lightKindFromInt(n int64) string {
	switch n {
	case 1:
		return "directional"
	case 2:
		return "point"
	case 3:
		return "spot"
	default:
		return "directional"
	}
}
