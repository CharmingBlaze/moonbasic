package vm

import (
	"sort"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

func normalizeHandleMethod(mn string) string {
	switch mn {
	case "SETPOSITION", "POSITION", "POS":
		return "SETPOS"
	case "SETROTATION", "ROTATION", "ROTATE", "ROT":
		return "SETROT"
	case "SETSCALE", "SCALE":
		return "SETSCALE"
	case "SETSIZE", "SIZE":
		return "SETSIZE"
	case "SETCOLOR", "COLOR", "COL":
		return "SETCOLOR"
	case "SETALPHA", "ALPHA", "A":
		return "SETALPHA"
	case "FOV":
		return "SETFOV"
	case "LOOK":
		return "LOOKAT"
	default:
		return mn
	}
}

// handleCallRegistryPrefix returns the dotted namespace prefix for handle-typed registry keys (e.g. CAMERA for Camera3D).
func handleCallRegistryPrefix(tag uint16) string {
	switch tag {
	case heap.TagCamera:
		return "CAMERA."
	case heap.TagCamera2D:
		return "CAMERA2D."
	case heap.TagRenderTexture:
		return "RENDERTARGET."
	case heap.TagInstancedModel:
		return "INSTANCE."
	case heap.TagMatrix:
		return "TRANSFORM."
	case heap.TagPhysicsBody:
		return "BODY3D."
	case heap.TagBody2D:
		return "PHYSICS2D."
	case heap.TagPeer:
		return "PEER."
	case heap.TagSound, heap.TagMusic, heap.TagAudioStream, heap.TagWave:
		return "AUDIO."
	case heap.TagTexture:
		return "TEXTURE."
	case heap.TagImage:
		return "IMAGE."
	case heap.TagMesh:
		return "MESH."
	case heap.TagFont:
		return "FONT."
	case heap.TagSprite:
		return "SPRITE."
	case heap.TagLight:
		return "LIGHT."
	case heap.TagLight2D:
		return "LIGHT2D."
	case heap.TagParticle:
		return "PARTICLE."
	case heap.TagTilemap:
		return "TILEMAP."
	case heap.TagAtlas:
		return "ATLAS."
	case heap.TagCharController:
		return "CHARACTERREF."
	case heap.TagShape:
		return "SHAPEREF."
	case heap.TagKinematicBody:
		return "KINEMATICREF."
	case heap.TagStaticBody:
		return "STATICREF."
	case heap.TagTriggerBody:
		return "TRIGGERREF."
	case heap.TagModel, heap.TagLODModel:
		return "MODEL."
	default:
		return ""
	}
}

// handleCallBuiltin maps heap TypeTag + script method to a registry command key and whether
// the receiver handle is passed as the first argument to that builtin.
func handleCallBuiltin(tag uint16, method string) (registryKey string, prependReceiver bool, ok bool) {
	mn := normalizeHandleMethod(strings.ToUpper(strings.TrimSpace(method)))
	switch tag {
	case heap.TagCamera:
		switch mn {
		case "END":
			return "CAMERA.END", false, true
		case "BEGIN", "SETPOS", "SETTARGET", "LOOKAT", "SETFOV", "SETPROJECTION", "MOVE", "GETRAY", "GETVIEWRAY", "GETMATRIX",
			"GETPOS", "GETTARGET", "SETUP", "FREE", "WORLDTOSCREEN", "ISONSCREEN", "MOUSERAY", "ZOOM", "ORBIT", "SETORBIT",
			"YAW", "GETYAW", "USEMOUSEORBIT", "USEORBITRIGHTMOUSE", "SETORBITKEYS", "SETORBITLIMITS", "SETORBITSPEED", "SETORBITKEYSPEED":
			return "CAMERA." + mn, true, true
		case "SETROT":
			return "CAMERA.ROTATE", true, true
		case "TURN":
			return "CAMERA.TURN", true, true
		case "FOLLOW":
			return "CAMERA.FOLLOWENTITY", true, true
		case "SHAKE":
			return "CAMERA.SHAKE", true, true
		case "PICK":
			return "CAMERA.PICK", true, true
		}
	case heap.TagEntityRef:
		switch mn {
		case "SETPOS":
			return "ENTITY.SETPOSITION", true, true
		case "MOVE":
			return "ENTITY.MOVE", true, true
		case "PUSH":
			return "ENTITY.PUSH", true, true
		case "JUMP":
			return "ENTITY.JUMP", true, true
		case "ISGROUNDED", "GROUNDED":
			return "ENTITY.GROUNDED", true, true
		case "SQUASH":
			return "ENTITY.SQUASH", true, true
		case "ADDPHYSICS":
			return "ENTITY.ADDPHYSICS", true, true
		case "SETBOUNCINESS", "BOUNCINESS":
			return "ENTITY.SETBOUNCINESS", true, true
		case "SETSCALE":
			return "ENTITY.SCALE", true, true
		case "SETROT":
			return "ENTITY.ROTATEENTITY", true, true
		case "TURN":
			return "ENTITY.TURNENTITY", true, true
		case "SETCOLOR":
			return "ENTITY.COLOR", true, true
		case "SETALPHA":
			return "ENTITY.ALPHA", true, true
		case "FREE":
			return "ENTITY.FREE", true, true
		case "HIDE":
			return "ENTITY.HIDE", true, true
		case "SHOW":
			return "ENTITY.SHOW", true, true
		case "MOVEWITHCAMERA":
			return "ENTITY.MOVEWITHCAMERA", true, true
		}
	case heap.TagCamera2D:
		switch mn {
		case "END":
			return "CAMERA2D.END", false, true
		case "BEGIN", "SETTARGET", "SETOFFSET", "SETZOOM", "SETROTATION", "GETMATRIX", "WORLDTOSCREEN", "SCREENTOWORLD", "FREE", "FOLLOW":
			return "CAMERA2D." + mn, true, true
		case "SETPOS", "TARGET":
			return "CAMERA2D.SETTARGET", true, true
		case "OFFSET":
			return "CAMERA2D.SETOFFSET", true, true
		case "ZOOM":
			return "CAMERA2D.SETZOOM", true, true
		case "SETROT", "ROT":
			return "CAMERA2D.SETROTATION", true, true
		}
	case heap.TagRenderTexture:
		switch mn {
		case "END":
			return "RENDERTARGET.END", false, true
		case "BEGIN", "FREE", "TEXTURE":
			return "RENDERTARGET." + mn, true, true
		}
	case heap.TagTilemap:
		switch mn {
		case "FREE", "SETTILESIZE", "DRAW", "DRAWLAYER", "GETTILE", "SETTILE", "ISSOLID", "ISSOLIDCATEGORY",
			"WIDTH", "HEIGHT", "LAYERCOUNT", "LAYERNAME", "COLLISIONAT", "SETCOLLISION", "MERGECOLLISIONLAYER":
			return "TILEMAP." + mn, true, true
		}
	case heap.TagAtlas:
		switch mn {
		case "FREE", "GETSPRITE":
			return "ATLAS." + mn, true, true
		}
	case heap.TagLight2D:
		switch mn {
		case "FREE", "SETPOS", "SETCOLOR", "SETRADIUS", "SETINTENSITY":
			return "LIGHT2D." + mn, true, true
		}
	case heap.TagPhysicsBody:
		switch mn {
		case "SETPOS":
			return "BODY3D.SETPOS", true, true
		case "SETROT":
			return "BODY3D.SETROT", true, true
		case "SETSCALE":
			return "BODY3D.SETSCALE", true, true
		case "SETVELOCITY", "SETVEL":
			return "BODY3D.SETLINEARVELOCITY", true, true
		case "ADDFORCE", "FORCE":
			return "BODY3D.ADDFORCE", true, true
		case "ADDIMPULSE", "IMPULSE":
			return "BODY3D.ADDIMPULSE", true, true
		case "FREE":
			return "BODY3D.FREE", true, true
		}
	case heap.TagBody2D:
		switch mn {
		case "SETPOS":
			return "PHYSICS2D.SETBODYPOSITION", true, true
		case "SETVELOCITY", "SETVEL":
			return "PHYSICS2D.SETBODYLINEARVELOCITY", true, true
		case "FREE":
			return "PHYSICS2D.DESTROYBODY", true, true
		}
	case heap.TagPeer:
		switch mn {
		case "SEND", "SENDPACKET":
			return "PEER.SENDPACKET", true, true
		case "PING":
			return "PEER.PING", true, true
		case "DISCONNECT":
			return "PEER.DISCONNECT", true, true
		case "IP":
			return "PEER.IP", true, true
		}
	case heap.TagCharController:
		switch mn {
		case "SETPOS", "POSITION":
			return "CHARACTERREF.SETPOSITION", true, true
		case "MOVE":
			return "CHARACTERREF.MOVE", true, true
		case "UPDATE":
			return "CHARACTERREF.UPDATE", true, true
		case "SETVELOCITY", "SETVEL":
			return "CHARACTERREF.SETVELOCITY", true, true
		case "ADDVELOCITY", "ADDVEL":
			return "CHARACTERREF.ADDVELOCITY", true, true
		case "SETMAXSLOPE", "SETSLOPE":
			return "CHARACTERREF.SETMAXSLOPE", true, true
		case "SETSTEPHEIGHT", "SETSTEP":
			return "CHARACTERREF.SETSTEPHEIGHT", true, true
		case "SETSNAPDISTANCE", "SNAP", "SETSTICKDOWN":
			return "CHARACTERREF.SETSNAPDISTANCE", true, true
		case "ISGROUNDED", "GROUNDED":
			return "CHARACTERREF.ISGROUNDED", true, true
		case "ONSLOPE":
			return "CHARACTERREF.ONSLOPE", true, true
		case "ONWALL":
			return "CHARACTERREF.ONWALL", true, true
		case "GETSLOPEANGLE", "SLOPEANGLE":
			return "CHARACTERREF.GETSLOPEANGLE", true, true
		case "JUMP":
			return "CHARACTERREF.JUMP", true, true
		case "GETPOSITION", "GETPOS":
			return "CHARACTERREF.GETPOSITION", true, true
		case "GETSPEED", "SPEED":
			return "CHARACTERREF.GETSPEED", true, true
		case "SETGRAVITY", "SETGRAVITYSCALE":
			return "CHARACTERREF.SETGRAVITY", true, true
		case "SETFRICTION":
			return "CHARACTERREF.SETFRICTION", true, true
		case "SETBOUNCE", "SETBOUNCINESS":
			return "CHARACTERREF.SETBOUNCE", true, true
		case "SETPADDING":
			return "CHARACTERREF.SETPADDING", true, true
		case "GETGROUNDSTATE":
			return "CHARACTERREF.GETGROUNDSTATE", true, true
		case "FREE":
			return "CHARACTERREF.FREE", true, true
		case "MOVEWITHCAMERA", "MOVEWITHCAM":
			return "CHARACTERREF.MOVEWITHCAMERA", true, true
		}
	case heap.TagShape:
		switch mn {
		case "FREE":
			return "SHAPEREF.FREE", true, true
		}
	case heap.TagKinematicBody, heap.TagStaticBody, heap.TagTriggerBody:
		switch mn {
		case "SETPOS":
			return "BODYREF.SETPOSITION", true, true
		case "SETROT":
			return "BODYREF.SETROTATION", true, true
		case "SETLAYER":
			return "BODYREF.SETLAYER", true, true
		case "ENABLECOLLISION":
			return "BODYREF.ENABLECOLLISION", true, true
		case "SETVELOCITY", "SETVEL":
			if tag == heap.TagKinematicBody {
				return "KINEMATICREF.SETVELOCITY", true, true
			}
		case "UPDATE":
			if tag == heap.TagKinematicBody {
				return "KINEMATICREF.UPDATE", true, true
			}
		case "FREE":
			return "BODYREF.FREE", true, true
		}
	case heap.TagSprite:
		switch mn {
		case "SETPOS", "DRAW", "DEFANIM", "PLAYANIM", "UPDATEANIM", "HIT", "FREE":
			return "SPRITE." + mn, true, true
		}
	case heap.TagLight:
		switch mn {
		case "SETDIR", "SETSHADOW", "FREE", "SETCOLOR", "SETINTENSITY", "SETPOSITION", "SETPOS",
			"SETTARGET", "SETSHADOWBIAS", "SETINNERCONE", "SETOUTERCONE", "SETRANGE", "ENABLE", "ISENABLED":
			return "LIGHT." + mn, true, true
		}
	case heap.TagModel, heap.TagLODModel:
		switch mn {
		case "SETPOS":
			return "MODEL.SETPOS", true, true
		case "SETROT":
			return "MODEL.SETROT", true, true
		case "SETSCALE":
			return "MODEL.SETSCALE", true, true
		case "DRAW":
			return "MODEL.DRAW", true, true
		}
	case heap.TagInstancedModel:
		switch mn {
		case "DRAW":
			return "MODEL.DRAW", true, true
		case "FREE":
			return "INSTANCE.FREE", true, true
		case "COUNT":
			return "INSTANCE.COUNT", true, true
		case "SETINSTANCEPOS", "SETPOS":
			return "INSTANCE.SETPOS", true, true
		case "SETINSTANCESCALE", "SETSCALE":
			return "INSTANCE.SETSCALE", true, true
		case "SETROT":
			return "INSTANCE.SETROT", true, true
		case "SETMATRIX":
			return "INSTANCE.SETMATRIX", true, true
		case "SETCOLOR":
			return "INSTANCE.SETCOLOR", true, true
		case "UPDATEINSTANCES", "UPDATEBUFFER":
			return "INSTANCE.UPDATEBUFFER", true, true
		case "SETCULLDISTANCE":
			return "INSTANCE.SETCULLDISTANCE", true, true
		case "DRAWLOD":
			return "INSTANCE.DRAWLOD", true, true
		}
	case heap.TagParticle:
		switch mn {
		case "SETTEXTURE":
			return "PARTICLE.SETTEXTURE", true, true
		case "SETEMITRATE":
			return "PARTICLE.SETEMITRATE", true, true
		case "SETRATE":
			return "PARTICLE.SETRATE", true, true
		case "SETLIFETIME":
			return "PARTICLE.SETLIFETIME", true, true
		case "SETVELOCITY":
			return "PARTICLE.SETVELOCITY", true, true
		case "SETDIRECTION":
			return "PARTICLE.SETDIRECTION", true, true
		case "SETSPREAD":
			return "PARTICLE.SETSPREAD", true, true
		case "SETSPEED":
			return "PARTICLE.SETSPEED", true, true
		case "SETSTARTSIZE":
			return "PARTICLE.SETSTARTSIZE", true, true
		case "SETENDSIZE":
			return "PARTICLE.SETENDSIZE", true, true
		case "SETCOLOR":
			return "PARTICLE.SETCOLOR", true, true
		case "SETSTARTCOLOR":
			return "PARTICLE.SETSTARTCOLOR", true, true
		case "SETCOLOREND":
			return "PARTICLE.SETCOLOREND", true, true
		case "SETENDCOLOR":
			return "PARTICLE.SETENDCOLOR", true, true
		case "SETSIZE":
			return "PARTICLE.SETSIZE", true, true
		case "SETGRAVITY":
			return "PARTICLE.SETGRAVITY", true, true
		case "SETPOS":
			return "PARTICLE.SETPOS", true, true
		case "SETBURST":
			return "PARTICLE.SETBURST", true, true
		case "SETBILLBOARD":
			return "PARTICLE.SETBILLBOARD", true, true
		case "PLAY":
			return "PARTICLE.PLAY", true, true
		case "STOP":
			return "PARTICLE.STOP", true, true
		case "UPDATE":
			return "PARTICLE.UPDATE", true, true
		case "DRAW":
			return "PARTICLE.DRAW", true, true
		case "ISALIVE":
			return "PARTICLE.ISALIVE", true, true
		case "COUNT":
			return "PARTICLE.COUNT", true, true
		case "FREE":
			return "PARTICLE.FREE", true, true
		}
	case heap.TagMatrix:
		if mn == "SETROTATION" {
			return "TRANSFORM.SETROTATION", true, true
		}
	case heap.TagMesh:
		switch mn {
		case "DRAW", "DRAWROTATED", "FREE":
			return "MESH." + mn, true, true
		case "VERTEXCOUNT", "TRIANGLECOUNT":
			return "MESH." + mn, true, true
		}
	case heap.TagDrawPrim3D:
		switch mn {
		case "POS", "SETPOS":
			return "DRAWPRIM3D.POS", true, true
		case "SIZE":
			return "DRAWPRIM3D.SIZE", true, true
		case "COLOR":
			return "DRAWPRIM3D.COLOR", true, true
		case "COL":
			return "DRAWPRIM3D.COL", true, true
		case "WIRE":
			return "DRAWPRIM3D.WIRE", true, true
		case "RADIUS":
			return "DRAWPRIM3D.RADIUS", true, true
		case "ENDPOINT":
			return "DRAWPRIM3D.ENDPOINT", true, true
		case "CYL":
			return "DRAWPRIM3D.CYL", true, true
		case "BBOX":
			return "DRAWPRIM3D.BBOX", true, true
		case "SLICES":
			return "DRAWPRIM3D.SLICES", true, true
		case "RINGS":
			return "DRAWPRIM3D.RINGS", true, true
		case "GRID":
			return "DRAWPRIM3D.GRID", true, true
		case "SETRAY":
			return "DRAWPRIM3D.SETRAY", true, true
		case "SETTEXTURE":
			return "DRAWPRIM3D.SETTEXTURE", true, true
		case "SRCTEX":
			return "DRAWPRIM3D.SRCTEX", true, true
		case "DRAW":
			return "DRAWPRIM3D.DRAW", true, true
		case "FREE":
			return "DRAWPRIM3D.FREE", true, true
		}
	case heap.TagDrawPrim2D:
		switch mn {
		case "POS", "SETPOS":
			return "DRAWPRIM2D.POS", true, true
		case "SIZE":
			return "DRAWPRIM2D.SIZE", true, true
		case "COLOR":
			return "DRAWPRIM2D.COLOR", true, true
		case "COL":
			return "DRAWPRIM2D.COL", true, true
		case "OUTLINE":
			return "DRAWPRIM2D.OUTLINE", true, true
		case "P2":
			return "DRAWPRIM2D.P2", true, true
		case "P3":
			return "DRAWPRIM2D.P3", true, true
		case "RING":
			return "DRAWPRIM2D.RING", true, true
		case "SEGS":
			return "DRAWPRIM2D.SEGS", true, true
		case "SIDES":
			return "DRAWPRIM2D.SIDES", true, true
		case "ROT":
			return "DRAWPRIM2D.ROT", true, true
		case "THICK":
			return "DRAWPRIM2D.THICK", true, true
		case "DRAW":
			return "DRAWPRIM2D.DRAW", true, true
		case "FREE":
			return "DRAWPRIM2D.FREE", true, true
		}
	case heap.TagTextDraw:
		switch mn {
		case "POS", "SETPOS":
			return "TEXTDRAW.POS", true, true
		case "SIZE":
			return "TEXTDRAW.SIZE", true, true
		case "COLOR":
			return "TEXTDRAW.COLOR", true, true
		case "COL":
			return "TEXTDRAW.COL", true, true
		case "SETTEXT":
			return "TEXTDRAW.SETTEXT", true, true
		case "DRAW":
			return "TEXTDRAW.DRAW", true, true
		case "FREE":
			return "TEXTDRAW.FREE", true, true
		}
	case heap.TagTextDrawEx:
		switch mn {
		case "POS", "SETPOS":
			return "TEXTEXOBJ.POS", true, true
		case "SIZE":
			return "TEXTEXOBJ.SIZE", true, true
		case "SPACING":
			return "TEXTEXOBJ.SPACING", true, true
		case "COLOR":
			return "TEXTEXOBJ.COLOR", true, true
		case "COL":
			return "TEXTEXOBJ.COLOR", true, true
		case "SETTEXT":
			return "TEXTEXOBJ.SETTEXT", true, true
		case "DRAW":
			return "TEXTEXOBJ.DRAW", true, true
		case "FREE":
			return "TEXTEXOBJ.FREE", true, true
		}
	case heap.TagTextureDraw:
		switch mn {
		case "POS", "SETPOS":
			return "DRAWTEX2.POS", true, true
		case "COLOR", "COL":
			return "DRAWTEX2.COLOR", true, true
		case "SETTEXTURE": // Covers DRAWTEX2, DRAWTEXREC, DRAWTEXPRO
			return "DRAWTEX2.SETTEXTURE", true, true
		case "DRAW":
			return "DRAWTEX2.DRAW", true, true
		case "FREE":
			return "DRAWTEX2.FREE", true, true
		case "SRC":
			return "DRAWTEXREC.SRC", true, true
		case "DST":
			return "DRAWTEXPRO.DST", true, true
		case "ORIGIN":
			return "DRAWTEXPRO.ORIGIN", true, true
		case "ROT":
			return "DRAWTEXPRO.ROT", true, true
		}
	case heap.TagInputFacade:
		switch mn {
		case "DX":
			return "MOUSE.DX", true, true
		case "DY":
			return "MOUSE.DY", true, true
		case "WHEEL":
			return "MOUSE.WHEEL", true, true
		case "DOWN": // Mouse, Key
			return "MOUSE.DOWN", true, true
		case "PRESSED":
			return "MOUSE.PRESSED", true, true
		case "RELEASED":
			return "MOUSE.RELEASED", true, true
		case "HIT":
			return "KEY.HIT", true, true
		case "UP":
			return "KEY.UP", true, true
		case "AXIS":
			return "GAMEPAD.AXIS", true, true
		case "BUTTON":
			return "GAMEPAD.BUTTON", true, true
		}
	case heap.TagSound, heap.TagAudioStream, heap.TagWave:
		switch mn {
		case "PLAY":
			return "AUDIO.PLAY", true, true
		case "STOP":
			return "AUDIO.STOP", true, true
		case "SETVOLUME", "VOLUME":
			return "AUDIO.SETSOUNDVOLUME", true, true
		case "FREE":
			return "AUDIO.FREESOUND", true, true
		}
	case heap.TagMusic:
		switch mn {
		case "PLAY":
			return "AUDIO.PLAY", true, true
		case "STOP":
			return "AUDIO.STOP", true, true
		case "SETVOLUME", "VOLUME":
			return "AUDIO.SETMUSICVOLUME", true, true
		case "FREE":
			return "AUDIO.FREEMUSIC", true, true
		}
	case heap.TagTexture, heap.TagImage:
		pre := "TEXTURE"
		if tag == heap.TagImage {
			pre = "IMAGE"
		}
		switch mn {
		case "DRAW":
			return pre + ".DRAW", true, true
		case "WIDTH":
			return pre + ".WIDTH", true, true
		case "HEIGHT":
			return pre + ".HEIGHT", true, true
		case "FREE":
			return pre + ".FREE", true, true
		}
	case heap.TagFont:
		switch mn {
		case "WIDTH", "TEXTWIDTH":
			return "FONT.TEXTWIDTH", true, true
		case "HEIGHT", "TEXTHEIGHT":
			return "FONT.TEXTHEIGHT", true, true
		case "FREE":
			return "FONT.FREE", true, true
		}
	case heap.TagMoverFacade:
		switch mn {
		case "MOVEXZ":
			return "MOVER.MOVEXZ", true, true
		case "MOVESTEPX":
			return "MOVER.MOVESTEPX", true, true
		case "MOVESTEPZ":
			return "MOVER.MOVESTEPZ", true, true
		case "LAND":
			return "MOVER.LAND", true, true
		case "MOVEREL":
			return "MOVER.MOVEREL", true, true
		case "FREE":
			return "MOVER.FREE", true, true
		}
	}
	return "", false, false
}

// HandleCallSuggestions lists common script-side method names for a handle type (error hints).
func HandleCallSuggestions(tag uint16) []string {
	var out []string
	switch tag {
	case heap.TagCamera:
		out = []string{"Begin", "End", "FOV", "Free", "GetMatrix", "GetPos", "GetRay", "GetTarget", "GetViewRay", "GetYaw", "IsOnScreen",
			"Look", "LookAt", "MouseRay", "Move", "Orbit", "Pos", "SetFOV", "SetOrbit", "SetPos", "SetPosition", "SetProjection", "SetTarget", "SetUp", "WorldToScreen", "Yaw", "Zoom"}
	case heap.TagEntityRef:
		out = []string{"A", "Col", "Color", "Free", "Hide", "Move", "MoveWithCamera", "Pos", "Rot", "Scale", "SetBounciness", "Show", "Turn"}
	case heap.TagCamera2D:
		out = []string{"Begin", "End", "Free", "GetMatrix", "ScreenToWorld", "SetOffset", "SetRotation", "SetTarget", "SetZoom", "WorldToScreen"}
	case heap.TagRenderTexture:
		out = []string{"Begin", "End", "Free", "Texture"}
	case heap.TagTilemap:
		out = []string{"CollisionAt", "Draw", "DrawLayer", "Free", "GetTile", "Height", "IsSolid", "IsSolidCategory",
			"LayerCount", "LayerName", "MergeCollisionLayer", "SetCollision", "SetTile", "SetTileSize", "Width"}
	case heap.TagAtlas:
		out = []string{"Free", "GetSprite"}
	case heap.TagLight2D:
		out = []string{"Free", "SetColor", "SetIntensity", "SetPos", "SetRadius"}
	case heap.TagPhysicsBody:
		out = []string{"AddForce", "AddImpulse", "Force", "Free", "Impulse", "Pos", "Rot", "Scale", "SetPos", "SetPosition", "SetRot", "SetVelocity", "Vel", "Velocity"}
	case heap.TagBody2D:
		out = []string{"Free", "Pos", "SetPos", "SetVel", "SetVelocity", "Vel", "Velocity"}
	case heap.TagPeer:
		out = []string{"Disconnect", "IP", "Ping", "Send", "SendPacket"}
	case heap.TagSound:
		out = []string{"Free", "Pause", "Play", "Resume", "SetVolume", "Stop", "Volume"}
	case heap.TagMusic:
		out = []string{"Free", "Pause", "Play", "Resume", "SetVolume", "Stop", "Volume"}
	case heap.TagCharController:
		out = []string{"AddVel", "AddVelocity", "Free", "GetPos", "GetPosition", "GetSlopeAngle", "GetSpeed", "Grounded", "IsGrounded", "Jump", "Move", "MoveWithCamera", "OnSlope", "OnWall", "Pos", "SetFriction", "SetGravity", "SetMaxSlope", "SetPos", "SetPosition", "SetSnapDistance", "SetStepHeight", "SetVelocity", "SlopeAngle", "Snap", "Speed", "Update"}
	case heap.TagShape:
		out = []string{"Free"}
	case heap.TagKinematicBody:
		out = []string{"EnableCollision", "Free", "Pos", "Rot", "SetLayer", "SetPos", "SetPosition", "SetRot", "SetRotation", "SetVel", "SetVelocity", "Update"}
	case heap.TagStaticBody, heap.TagTriggerBody:
		out = []string{"EnableCollision", "Free", "Pos", "Rot", "SetLayer", "SetPos", "SetPosition", "SetRot", "SetRotation"}
	case heap.TagSprite:
		out = []string{"Draw", "Free", "SetPos", "SetPosition", "DefAnim", "PlayAnim", "UpdateAnim", "Hit"}
	case heap.TagModel, heap.TagLODModel:
		out = []string{"Draw", "SetPos", "SetPosition"}
	case heap.TagParticle:
		out = []string{"Count", "Draw", "Free", "IsAlive", "Play", "SetBillboard", "SetBurst", "SetColor", "SetColorEnd",
			"SetDirection", "SetEmitRate", "SetEndColor", "SetEndSize", "SetGravity", "SetLifetime", "SetPos", "SetPosition",
			"SetRate", "SetSize", "SetSpeed", "SetSpread", "SetStartColor", "SetStartSize", "SetTexture", "SetVelocity", "Stop", "Update"}
	case heap.TagInstancedModel:
		out = []string{"Count", "Draw", "DrawLOD", "Free", "SetColor", "SetCullDistance", "SetInstancePos", "SetInstanceScale",
			"SetMatrix", "SetPos", "SetRot", "SetScale", "UpdateBuffer", "UpdateInstances"}
	case heap.TagLight:
		out = []string{"Enable", "Free", "IsEnabled", "SetColor", "SetDir", "SetInnerCone", "SetIntensity",
			"SetOuterCone", "SetPos", "SetPosition", "SetRange", "SetShadow", "SetShadowBias", "SetTarget"}
	case heap.TagMatrix:
		out = []string{"SetRotation"}
	case heap.TagMesh:
		out = []string{"Draw", "DrawRotated", "Free", "TriangleCount", "VertexCount"}
	case heap.TagFont:
		out = []string{"Free", "Height", "TextHeight", "TextWidth", "Width"}
	case heap.TagDrawPrim3D:
		out = []string{"Pos", "Size", "Color", "Col", "Wire", "Radius", "EndPoint", "Cyl", "BBox", "Slices", "Rings", "Grid", "SetRay", "SetTexture", "SrcTex", "Draw", "Free"}
	case heap.TagDrawPrim2D:
		out = []string{"Pos", "Size", "Color", "Col", "Outline", "P2", "P3", "Draw", "Free"}
	case heap.TagTextDraw:
		out = []string{"Pos", "Size", "Color", "Col", "SetText", "Draw", "Free"}
	case heap.TagTextDrawEx:
		out = []string{"Pos", "Size", "Spacing", "Color", "Col", "SetText", "Draw", "Free"}
	case heap.TagTextureDraw:
		out = []string{"Pos", "Color", "Col", "SetTexture", "Src", "Dst", "Origin", "Rot", "Draw", "Free"}
	case heap.TagInputFacade:
		out = []string{"DX", "DY", "Wheel", "Down", "Pressed", "Released", "Hit", "Up", "Axis", "Button"}
	case heap.TagMoverFacade:
		out = []string{"MoveXZ", "MoveStepX", "MoveStepZ", "Land", "MoveRel", "Free"}
	default:
		return nil
	}
	sort.Strings(out)
	return out
}

func filterRegistryKeysByPrefix(keys []string, prefix string) []string {
	pu := strings.ToUpper(prefix)
	var out []string
	for _, k := range keys {
		if strings.HasPrefix(strings.ToUpper(k), pu) {
			out = append(out, k)
		}
	}
	return out
}

// formatHandleCallError enriches a failed handle method dispatch with type-specific hints.
func (v *VM) formatHandleCallError(tag uint16, typeName, methodName, callKey string, mapped bool, err error) string {
	msg := err.Error()
	if mapped {
		return msg
	}
	prefix := handleCallRegistryPrefix(tag)
	keys := v.Registry.CommandKeys()
	prefixed := filterRegistryKeysByPrefix(keys, prefix)
	if alt, ok := runtime.BestSimilarCommand(callKey, prefixed, 3); ok {
		return msg + "\n  Did you mean " + alt + "?"
	}
	if sug := HandleCallSuggestions(tag); len(sug) > 0 {
		return msg + "\n  Hint: For this handle type use methods like " + strings.Join(sug, ", ") + "."
	}
	return msg + "\n  Hint: See docs/API_CONSISTENCY.md for handle methods vs NS.COMMAND calls."
}
