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
	case "FREE", "DESTROY", "DELETE":
		return "FREE"
	// Model render / PBR / texture-stage short names (see MODEL.* in runtime/mbmodel3d)
	case "CULL":
		return "SETCULL"
	case "WIREFRAME":
		return "SETWIREFRAME"
	case "LIGHTING":
		return "SETLIGHTING"
	case "FOG":
		return "SETFOG"
	case "BLEND":
		return "SETBLEND"
	case "DEPTH":
		return "SETDEPTH"
	case "DIFFUSE":
		return "SETDIFFUSE"
	case "SPECULAR":
		return "SETSPECULAR"
	case "SPECULARPOW":
		return "SETSPECULARPOW"
	case "EMISSIVE":
		return "SETEMISSIVE"
	case "AMBIENT", "AMBIENTCOLOR":
		return "SETAMBIENTCOLOR"
	case "METAL":
		return "SETMETAL"
	case "ROUGH":
		return "SETROUGH"
	case "TEXTURESTAGE":
		return "SETTEXTURESTAGE"
	case "STAGEBLEND":
		return "SETSTAGEBLEND"
	case "STAGESCROLL":
		return "SETSTAGESCROLL"
	case "STAGESCALE":
		return "SETSTAGESCALE"
	case "STAGEROTATE":
		return "SETSTAGEROTATE"
	case "LODDISTANCES":
		return "SETLODDISTANCES"
	case "GPUSKINNING":
		return "SETGPUSKINNING"
	case "MATERIALTEXTURE":
		return "SETMATERIALTEXTURE"
	case "MATERIALSHADER":
		return "SETMATERIALSHADER"
	case "MODELMESHMATERIAL":
		return "SETMODELMESHMATERIAL"
	case "MATRIX":
		return "SETMATRIX"
	case "ATTACH":
		return "ATTACHTO"
	case "DETAIL":
		return "SETDETAIL"
	case "ASYNC", "ASYNCMESH":
		return "SETASYNCMESHBUILD"
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
	case heap.TagPhysicsBody, heap.TagPhysicsBuilder:
		return "BODY3D."
	case heap.TagBody2D:
		return "PHYSICS2D."
	case heap.TagPeer:
		return "PEER."
	case heap.TagSound, heap.TagMusic:
		return "AUDIO."
	case heap.TagAudioStream:
		return "AUDIOSTREAM."
	case heap.TagWave:
		return "WAVE."
	case heap.TagTexture:
		return "TEXTURE."
	case heap.TagImage, heap.TagImageSequence:
		return "IMAGE."
	case heap.TagMesh:
		return "MESH."
	case heap.TagMaterial:
		return "MATERIAL."
	case heap.TagShader:
		return "SHADER."
	case heap.TagSky:
		return "SKY."
	case heap.TagCloud:
		return "CLOUD."
	case heap.TagWeather:
		return "WEATHER."
	case heap.TagDecal:
		return "DECAL."
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
	case heap.TagWater:
		return "WATER."
	case heap.TagTerrain:
		return "TERRAIN."
	case heap.TagTilemap:
		return "TILEMAP."
	case heap.TagAtlas:
		return "ATLAS."
	case heap.TagCharController:
		return "CHARACTERREF."
	case heap.TagEntityRef:
		return "ENTITY."
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
	case heap.TagFile:
		return "FILE."
	case heap.TagArray:
		return "ARRAY."
	case heap.TagNav:
		return "NAV."
	case heap.TagNavAgent:
		return "NAVAGENT."
	case heap.TagPath:
		return "PATH."
	case heap.TagScatterSet:
		return "SCATTER."
	case heap.TagProp:
		return "PROP."
	case heap.TagTween:
		return "TWEEN."
	case heap.TagBiome:
		return "BIOME."
	case heap.TagNoise:
		return "NOISE."
	case heap.TagTable:
		return "TABLE."
	case heap.TagPool:
		return "POOL."
	case heap.TagJSON:
		return "JSON."
	case heap.TagCSV:
		return "CSV."
	case heap.TagDB, heap.TagDBStmt, heap.TagDBTx:
		return "DB."
	case heap.TagDBRows:
		return "ROWS."
	case heap.TagRng:
		return "RAND."
	case heap.TagMem:
		return "MEM."
	case heap.TagLobby:
		return "LOBBY."
	case heap.TagNetPacket:
		return "PACKET."
	case heap.TagHost:
		return "NET."
	case heap.TagEvent, heap.TagAutomationList:
		return "EVENT."
	case heap.TagSpriteGroup:
		return "SPRITEGROUP."
	case heap.TagSpriteLayer:
		return "SPRITELAYER."
	case heap.TagSpriteBatch:
		return "SPRITEBATCH."
	case heap.TagSpriteUI:
		return "SPRITEUI."
	case heap.TagParticle2D:
		return "PARTICLE2D."
	case heap.TagQuaternion:
		return "QUAT."
	case heap.TagColor:
		return "COLOR."
	case heap.TagVec2:
		return "VEC2."
	case heap.TagVec3:
		return "VEC3."
	case heap.TagPlayer2D:
		return "PLAYER2D."
	case heap.TagGameTimer, heap.TagGameTimerSim:
		return "TIMER."
	case heap.TagGameStopwatch:
		return "STOPWATCH."
	case heap.TagBTree:
		return "BTREE."
	case heap.TagSteerGroup:
		return "STEER."
	case heap.TagComputeShader, heap.TagShaderBuffer:
		return "COMPUTESHADER."
	case heap.TagJoint2D:
		return "JOINT2D."
	case heap.TagDrawPrim3D:
		return "DRAWPRIM3D."
	case heap.TagDrawPrim2D:
		return "DRAWPRIM2D."
	case heap.TagTextDraw:
		return "TEXTDRAW."
	case heap.TagTextDrawEx:
		return "TEXTEXOBJ."
	case heap.TagTextureDraw:
		return "DRAWTEX"
	case heap.TagMoverFacade:
		return "MOVER."
	case heap.TagTacticalGrid:
		return "GRID."
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
			"GETPOS", "GETROT", "GETTARGET", "GETFOV", "GETUP", "GETPROJECTION", "SETUP", "FREE", "WORLDTOSCREEN", "ISONSCREEN", "MOUSERAY", "ZOOM", "ORBIT", "SETORBIT",
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
		case "SETPOSTION":
			return "CAMERA.SETPOS", true, true
		case "LERPTO":
			return "CAMERA.LERPTO", true, true
		}
	case heap.TagEntityRef:
		switch mn {
		case "SETPOS":
			return "ENTITY.SETPOS", true, true
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
		case "SETCOLLISIONMESH", "COLLISIONMESH":
			return "ENTITY.SETSTATIC", true, true
		case "SETSTATIC", "STATIC":
			return "ENTITY.SETSTATIC", true, true
		case "SETDAMPING":
			return "BODY3D.SETDAMPING", true, true
		case "LOCKAXIS":
			return "BODY3D.LOCKAXIS", true, true
		case "SETGRAVITYFACTOR":
			return "BODY3D.SETGRAVITYFACTOR", true, true
		case "SETCCD":
			return "BODY3D.SETCCD", true, true
		case "SETSTEERING", "STEER":
			return "VEHICLE.SETSTEER", true, true
		case "SETTHROTTLE", "THROTTLE":
			return "VEHICLE.SETTHROTTLE", true, true
		case "CREATEVEHICLE":
			return "VEHICLE.CREATE", true, true
		case "SETLIFT":
			return "AERO.SETLIFT", true, true
		case "SETTHRUST":
			return "AERO.SETTHRUST", true, true
		case "CREATECHARACTER":
			return "CHARACTER.CREATE", true, true
		}
	case heap.TagCamera2D:
		switch mn {
		case "END":
			return "CAMERA2D.END", false, true
		case "BEGIN", "SETTARGET", "SETOFFSET", "SETZOOM", "SETROTATION", "GETPOS", "GETROTATION", "GETMATRIX", "GETZOOM", "GETOFFSET", "WORLDTOSCREEN", "SCREENTOWORLD", "FREE", "FOLLOW":
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
			"LAYERCOUNT", "LAYERNAME", "COLLISIONAT", "SETCOLLISION", "MERGECOLLISIONLAYER":
			return "TILEMAP." + mn, true, true
		case "WIDTH", "GETWIDTH":
			return "TILEMAP.WIDTH", true, true
		case "HEIGHT", "GETHEIGHT":
			return "TILEMAP.HEIGHT", true, true
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
		case "SETPOS", "POS":
			return "BODY3D.SETPOS", true, true
		case "GETPOS":
			return "BODY3D.GETPOS", true, true
		case "SETROT", "ROT":
			return "BODY3D.SETROT", true, true
		case "GETROT":
			return "BODY3D.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "BODY3D.SETSCALE", true, true
		case "GETSCALE":
			return "BODY3D.GETSCALE", true, true
		case "SETVELOCITY", "SETVEL", "VEL", "VELOCITY":
			return "BODY3D.SETVELOCITY", true, true
		case "GETVELOCITY", "GETVEL":
			return "BODY3D.GETVELOCITY", true, true
		case "SETANGULARVELOCITY", "SETANGULARVEL", "ANGULARVELOCITY", "ANGULARVEL", "ANGVEL":
			return "BODY3D.SETANGULARVEL", true, true
		case "GETANGULARVEL":
			return "BODY3D.GETANGULARVEL", true, true
		case "SETMASS", "MASS":
			return "BODY3D.SETMASS", true, true
		case "GETMASS":
			return "BODY3D.GETMASS", true, true
		case "SETLIFT":
			return "AERO.SETLIFT", true, true
		case "SETTHRUST":
			return "AERO.SETTHRUST", true, true
		case "SETDRAG":
			return "AERO.SETDRAG", true, true
		case "SETTORQUE", "APPLYTORQUE":
			return "BODY3D.APPLYTORQUE", true, true
		case "ADDFORCE", "FORCE", "APPLYFORCE":
			return "BODY3D.APPLYFORCE", true, true
		case "ADDIMPULSE", "IMPULSE", "APPLYIMPULSE":
			return "BODY3D.APPLYIMPULSE", true, true
		case "ACTIVATE":
			return "BODY3D.ACTIVATE", true, true
		case "DEACTIVATE":
			return "BODY3D.DEACTIVATE", true, true
		case "COLLIDED":
			return "BODY3D.COLLIDED", true, true
		case "COLLISIONOTHER":
			return "BODY3D.COLLISIONOTHER", true, true
		case "COLLISIONPOINT":
			return "BODY3D.COLLISIONPOINT", true, true
		case "COLLISIONNORMAL":
			return "BODY3D.COLLISIONNORMAL", true, true
		case "BUFFERINDEX":
			return "BODY3D.BUFFERINDEX", true, true
		case "SETLINEARVEL", "SETLINEARVELOCITY":
			return "BODY3D.SETLINEARVEL", true, true
		case "GETLINEARVEL", "GETLINEARVELOCITY":
			return "BODY3D.GETLINEARVEL", true, true
		case "SETGRAVITYFACTOR", "GRAVITYFACTOR":
			return "BODY3D.SETGRAVITYFACTOR", true, true
		case "SETDAMPING", "DAMPING":
			return "BODY3D.SETDAMPING", true, true
		case "SETFRICTION", "FRICTION":
			return "BODY3D.SETFRICTION", true, true
		case "SETCCD", "CCD":
			return "BODY3D.SETCCD", true, true
		case "SETRESTITUTION", "SETBOUNCE", "RESTITUTION", "BOUNCE":
			return "BODY3D.SETRESTITUTION", true, true
		case "GETFRICTION":
			return "BODY3D.GETFRICTION", true, true
		case "GETRESTITUTION", "GETBOUNCE":
			return "BODY3D.GETRESTITUTION", true, true
		case "GETGRAVITYFACTOR":
			return "BODY3D.GETGRAVITYFACTOR", true, true
		case "GETDAMPING":
			return "BODY3D.GETDAMPING", true, true
		case "GETCCD":
			return "BODY3D.GETCCD", true, true
		case "FREE":
			return "BODY3D.FREE", true, true
		}
	case heap.TagPhysicsBuilder:
		switch mn {
		case "ADDBOX":
			return "BODY3D.ADDBOX", true, true
		case "ADDSPHERE":
			return "BODY3D.ADDSPHERE", true, true
		case "ADDCAPSULE":
			return "BODY3D.ADDCAPSULE", true, true
		case "ADDMESH":
			return "BODY3D.ADDMESH", true, true
		case "COMMIT":
			return "BODY3D.COMMIT", true, true
		case "FREE":
			return "BODY3D.FREE", true, true
		}
	case heap.TagBody2D:
		switch mn {
		case "SETPOS", "POS":
			return "BODY2D.SETPOS", true, true
		case "GETPOS":
			return "BODY2D.GETPOS", true, true
		case "SETROT", "ROT":
			return "BODY2D.SETROT", true, true
		case "GETROT":
			return "BODY2D.GETROT", true, true
		case "SETVELOCITY", "SETVEL", "VEL", "VELOCITY", "LINEARVEL":
			return "BODY2D.SETLINEARVELOCITY", true, true
		case "GETVELOCITY", "GETVEL":
			return "BODY2D.GETLINEARVELOCITY", true, true
		case "SETANGULARVELOCITY", "SETANGULARVEL", "ANGULARVELOCITY", "ANGULARVEL", "ANGVEL":
			return "BODY2D.SETANGULARVELOCITY", true, true
		case "GETANGULARVELOCITY", "GETANGULARVEL":
			return "BODY2D.GETANGULARVELOCITY", true, true
		case "SETFRICTION", "FRICTION":
			return "BODY2D.SETFRICTION", true, true
		case "GETFRICTION":
			return "BODY2D.GETFRICTION", true, true
		case "SETRESTITUTION", "RESTITUTION", "BOUNCE":
			return "BODY2D.SETRESTITUTION", true, true
		case "GETRESTITUTION", "GETBOUNCE":
			return "BODY2D.GETRESTITUTION", true, true
		case "SETMASS", "MASS":
			return "BODY2D.SETMASS", true, true
		case "GETMASS":
			return "BODY2D.GETMASS", true, true
		case "FREE":
			return "BODY2D.FREE", true, true
		case "APPLYFORCE", "ADDFORCE", "FORCE":
			return "BODY2D.APPLYFORCE", true, true
		case "APPLYIMPULSE", "ADDIMPULSE", "IMPULSE":
			return "BODY2D.APPLYIMPULSE", true, true
		case "ADDCIRCLE":
			return "BODY2D.ADDCIRCLE", true, true
		case "ADDRECT":
			return "BODY2D.ADDRECT", true, true
		case "ADDPOLYGON":
			return "BODY2D.ADDPOLYGON", true, true
		case "COMMIT":
			return "BODY2D.COMMIT", true, true
		case "COLLIDED":
			return "BODY2D.COLLIDED", true, true
		case "COLLISIONOTHER":
			return "BODY2D.COLLISIONOTHER", true, true
		case "COLLISIONPOINT":
			return "BODY2D.COLLISIONPOINT", true, true
		case "COLLISIONNORMAL":
			return "BODY2D.COLLISIONNORMAL", true, true
		}
	case heap.TagJoint2D:
		switch mn {
		case "FREE":
			return "JOINT2D.FREE", true, true
		}
	case heap.TagPeer:
		switch mn {
		case "SEND":
			return "PEER.SEND", true, true
		case "SENDPACKET":
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
			return "CHARACTERREF.SETPOS", true, true
		case "MOVE":
			return "CHARACTERREF.MOVE", true, true
		case "UPDATE":
			return "CHARACTERREF.UPDATE", true, true
		case "SETVELOCITY", "SETVEL", "VEL", "VELOCITY":
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
		case "GETSLOPEANGLE":
			return "CHARACTERREF.GETSLOPEANGLE", true, true
		case "JUMP":
			return "CHARACTERREF.JUMP", true, true
		case "GETPOSITION", "GETPOS":
			return "CHARACTERREF.GETPOSITION", true, true
		case "GETSPEED":
			return "CHARACTERREF.GETSPEED", true, true
		case "SETGRAVITY", "SETGRAVITYSCALE":
			return "CHARACTERREF.SETGRAVITY", true, true
		case "SETFRICTION":
			return "CHARACTERREF.SETFRICTION", true, true
		case "SETBOUNCE", "SETBOUNCINESS":
			return "CHARACTERREF.SETBOUNCINESS", true, true
		case "GETBOUNCE", "GETBOUNCINESS":
			return "CHARACTERREF.GETBOUNCINESS", true, true
		case "SETPADDING":
			return "CHARACTERREF.SETPADDING", true, true
		case "GETGROUNDSTATE":
			return "CHARACTERREF.GETGROUNDSTATE", true, true
		case "GETMAXSLOPE", "GETSLOPE":
			return "CHARACTERREF.GETMAXSLOPE", true, true
		case "GETSTEPHEIGHT", "GETSTEP":
			return "CHARACTERREF.GETSTEPHEIGHT", true, true
		case "GETSNAPDISTANCE", "GETSNAP":
			return "CHARACTERREF.GETSNAPDISTANCE", true, true
		case "GETGRAVITY", "GETGRAVITYSCALE":
			return "CHARACTERREF.GETGRAVITY", true, true
		case "GETFRICTION":
			return "CHARACTERREF.GETFRICTION", true, true
		case "GETPADDING":
			return "CHARACTERREF.GETPADDING", true, true
		case "FREE":
			return "CHARACTERREF.FREE", true, true
		case "MOVEWITHCAMERA", "MOVEWITHCAM":
			return "CHARACTERREF.MOVEWITHCAMERA", true, true
		case "SETJUMPBUFFER", "JUMPBUFFER":
			return "CHARACTERREF.SETJUMPBUFFER", true, true
		case "SETAIRCONTROL", "AIRCONTROL":
			return "CHARACTERREF.SETAIRCONTROL", true, true
		case "SETGROUNDCONTROL", "GROUNDCONTROL":
			return "CHARACTERREF.SETGROUNDCONTROL", true, true
		case "SETLINEARVELOCITY", "SETLINEARVEL":
			return "CHARACTERREF.SETLINEARVELOCITY", true, true
		case "UPDATEMOVE":
			return "CHARACTERREF.UPDATEMOVE", true, true
		case "DRAINCONTACTS":
			return "CHARACTERREF.DRAINCONTACTS", true, true
		case "SETCONTACTLISTENER":
			return "CHARACTERREF.SETCONTACTLISTENER", true, true
		case "SETSETTING":
			return "CHARACTERREF.SETSETTING", true, true
		case "GETCEILING":
			return "CHARACTERREF.GETCEILING", true, true
		case "GETISSLIDING":
			return "CHARACTERREF.GETISSLIDING", true, true
		case "GETGROUNDVELOCITY":
			return "CHARACTERREF.GETGROUNDVELOCITY", true, true
		}
	case heap.TagTween:
		switch mn {
		case "TO":
			return "TWEEN.TO", true, true
		case "THEN":
			return "TWEEN.THEN", true, true
		case "ONCOMPLETE":
			return "TWEEN.ONCOMPLETE", true, true
		case "START":
			return "TWEEN.START", true, true
		case "STOP":
			return "TWEEN.STOP", true, true
		case "LOOP":
			return "TWEEN.LOOP", true, true
		case "GETLOOP":
			return "TWEEN.GETLOOP", true, true
		case "YOYO":
			return "TWEEN.YOYO", true, true
		case "GETYOYO":
			return "TWEEN.GETYOYO", true, true
		case "UPDATE":
			return "TWEEN.UPDATE", true, true
		case "FREE":
			return "TWEEN.FREE", true, true
		}
	case heap.TagRay:
		switch mn {
		case "SETPOS", "POS", "POSITION":
			return "RAY.SETPOS", true, true
		case "SETDIR", "DIR":
			return "RAY.SETDIR", true, true
		case "GETPOS":
			return "RAY.GETPOS", true, true
		case "GETDIR":
			return "RAY.GETDIR", true, true
		case "HITSPHERE", "CHECKSPHERE":
			return "RAY.HITSPHERE", true, true
		case "HITBOX", "CHECKBOX":
			return "RAY.HITBOX", true, true
		case "HITPLANE", "CHECKPLANE":
			return "RAY.HITPLANE", true, true
		case "HITTRIANGLE", "CHECKTRIANGLE":
			return "RAY.HITTRIANGLE", true, true
		case "HITMESH", "CHECKMESH":
			return "RAY.HITMESH", true, true
		case "HITMODEL", "CHECKMODEL":
			return "RAY.HITMODEL", true, true
		case "INTERSECTSMODEL":
			return "RAY.INTERSECTSMODEL", true, true
		case "FREE":
			return "RAY.FREE", true, true
		}
	case heap.TagBBox:
		switch mn {
		case "SETMIN", "MIN":
			return "BBOX.SETMIN", true, true
		case "SETMAX", "MAX":
			return "BBOX.SETMAX", true, true
		case "GETMIN":
			return "BBOX.GETMIN", true, true
		case "GETMAX":
			return "BBOX.GETMAX", true, true
		case "FREE":
			return "BBOX.FREE", true, true
		case "CHECK":
			return "BBOX.CHECK", true, true
		case "CHECKSPHERE":
			return "BBOX.CHECKSPHERE", true, true
		}
	case heap.TagBSphere:
		switch mn {
		case "SETPOS", "POS", "POSITION":
			return "BSPHERE.SETPOS", true, true
		case "SETRADIUS", "RADIUS":
			return "BSPHERE.SETRADIUS", true, true
		case "GETPOS":
			return "BSPHERE.GETPOS", true, true
		case "GETRADIUS":
			return "BSPHERE.GETRADIUS", true, true
		case "FREE":
			return "BSPHERE.FREE", true, true
		case "CHECK":
			return "BSPHERE.CHECK", true, true
		case "CHECKBOX":
			return "BSPHERE.CHECKBOX", true, true
		}
	case heap.TagShape:
		switch mn {
		case "FREE":
			return "SHAPEREF.FREE", true, true
		}
	case heap.TagKinematicBody, heap.TagStaticBody, heap.TagTriggerBody:
		switch mn {
		case "SETPOS", "POSITION":
			return "BODYREF.SETPOS", true, true
		case "SETROT":
			return "BODYREF.SETROTATION", true, true
		case "SETLAYER":
			return "BODYREF.SETLAYER", true, true
		case "ENABLECOLLISION":
			return "BODYREF.ENABLECOLLISION", true, true
		case "SETVELOCITY", "SETVEL", "VEL", "VELOCITY":
			if tag == heap.TagKinematicBody {
				return "KINEMATICREF.SETVELOCITY", true, true
			}
		case "GETVELOCITY", "GETVEL":
			if tag == heap.TagKinematicBody {
				return "KINEMATICREF.GETVELOCITY", true, true
			}
		case "UPDATE":
			if tag == heap.TagKinematicBody {
				return "KINEMATICREF.UPDATE", true, true
			}
		case "FREE":
			return "BODYREF.FREE", true, true
		}
	case heap.TagNav:
		switch mn {
		case "SETGRID":
			return "NAV.SETGRID", true, true
		case "ADDTERRAIN":
			return "NAV.ADDTERRAIN", true, true
		case "ADDOBSTACLE":
			return "NAV.ADDOBSTACLE", true, true
		case "BUILD":
			return "NAV.BUILD", true, true
		case "DEBUGDRAW":
			return "NAV.DEBUGDRAW", true, true
		case "FINDPATH":
			return "NAV.FINDPATH", true, true
		case "FREE":
			return "NAV.FREE", true, true
		}
	case heap.TagPath:
		switch mn {
		case "ISVALID":
			return "PATH.ISVALID", true, true
		case "NODECOUNT", "SETSIZE":
			return "PATH.NODECOUNT", true, true
		case "NODEX":
			return "PATH.NODEX", true, true
		case "NODEY":
			return "PATH.NODEY", true, true
		case "NODEZ":
			return "PATH.NODEZ", true, true
		case "FREE":
			return "PATH.FREE", true, true
		}
	case heap.TagNavAgent:
		switch mn {
		case "SETPOS", "POSITION":
			return "NAVAGENT.SETPOS", true, true
		case "GETPOS":
			return "NAVAGENT.GETPOS", true, true
		case "SETROT", "ROT":
			return "NAVAGENT.SETROT", true, true
		case "GETROT":
			return "NAVAGENT.GETROT", true, true
		case "SETSPEED", "SPEED":
			return "NAVAGENT.SETSPEED", true, true
		case "GETSPEED":
			return "NAVAGENT.GETSPEED", true, true
		case "SETMAXFORCE", "MAXFORCE":
			return "NAVAGENT.SETMAXFORCE", true, true
		case "GETMAXFORCE":
			return "NAVAGENT.GETMAXFORCE", true, true
		case "APPLYFORCE":
			return "NAVAGENT.APPLYFORCE", true, true
		case "MOVETO":
			return "NAVAGENT.MOVETO", true, true
		case "UPDATE":
			return "NAVAGENT.UPDATE", true, true
		case "ISATDESTINATION":
			return "NAVAGENT.ISATDESTINATION", true, true
		case "X":
			return "NAVAGENT.X", true, true
		case "Y":
			return "NAVAGENT.Y", true, true
		case "Z":
			return "NAVAGENT.Z", true, true
		case "FREE":
			return "NAVAGENT.FREE", true, true
		case "STOP":
			return "NAVAGENT.STOP", true, true
		}
	case heap.TagBTree:
		switch mn {
		case "SEQUENCE", "SEQ":
			return "BTREE.SEQUENCE", true, true
		case "ADDCONDITION", "COND":
			return "BTREE.ADDCONDITION", true, true
		case "ADDACTION", "ACT":
			return "BTREE.ADDACTION", true, true
		case "RUN":
			return "BTREE.RUN", true, true
		case "FREE":
			return "BTREE.FREE", true, true
		}
	case heap.TagSteerGroup:
		switch mn {
		case "GROUPADD", "ADD":
			return "STEER.GROUPADD", true, true
		case "GROUPCLEAR", "CLEAR":
			return "STEER.GROUPCLEAR", true, true
		}
	case heap.TagSprite:
		switch mn {
		case "SETPOS", "POS":
			return "SPRITE.SETPOS", true, true
		case "GETPOS":
			return "SPRITE.GETPOS", true, true
		case "SETSCALE", "SCALE":
			return "SPRITE.SETSCALE", true, true
		case "GETSCALE":
			return "SPRITE.GETSCALE", true, true
		case "SETROT", "ROT":
			return "SPRITE.SETROT", true, true
		case "GETROT":
			return "SPRITE.GETROT", true, true
		case "SETCOLOR", "SETCOL", "COL":
			return "SPRITE.SETCOLOR", true, true
		case "GETCOLOR", "GETCOL":
			return "SPRITE.GETCOLOR", true, true
		case "SETALPHA", "ALPHA":
			return "SPRITE.SETALPHA", true, true
		case "GETALPHA":
			return "SPRITE.GETALPHA", true, true
		case "DRAW", "DEFANIM", "PLAYANIM", "UPDATEANIM", "SETFRAME", "SETORIGIN", "PLAY", "POINTHIT", "FREE":
			return "SPRITE." + mn, true, true
		case "HIT", "COLLIDE":
			return "SPRITE.HIT", true, true
		}
	case heap.TagLight:
		switch mn {
		case "SETPOS", "SETPOSITION", "POS", "POSITION":
			return "LIGHT.SETPOS", true, true
		case "GETPOS":
			return "LIGHT.GETPOS", true, true
		case "SETDIR", "DIR", "ROT", "SETROT":
			return "LIGHT.SETDIR", true, true
		case "GETDIR":
			return "LIGHT.GETDIR", true, true
		case "SETCOLOR", "SETCOL", "COL", "COLOR":
			return "LIGHT.SETCOLOR", true, true
		case "GETCOLOR", "GETCOL":
			return "LIGHT.GETCOLOR", true, true
		case "SETINTENSITY", "INTENSITY", "ENERGY", "SETENERGY":
			return "LIGHT.SETINTENSITY", true, true
		case "GETINTENSITY":
			return "LIGHT.GETINTENSITY", true, true
		case "SETRANGE", "RANGE":
			return "LIGHT.SETRANGE", true, true
		case "SETSHADOW", "SHADOW":
			return "LIGHT.SETSHADOW", true, true
		case "SETTARGET", "SETSHADOWBIAS", "SETINNERCONE", "SETOUTERCONE", "ENABLE", "SETSTATE", "ISENABLED":
			return "LIGHT." + mn, true, true
		case "GETSHADOW", "GETRANGE", "GETINNERCONE", "GETOUTERCONE", "GETENERGY":
			return "LIGHT." + mn, true, true
		case "FREE":
			return "LIGHT.FREE", true, true
		}
	case heap.TagModel, heap.TagLODModel:
		switch mn {
		case "SETPOS", "POS":
			return "MODEL.SETPOS", true, true
		case "GETPOS":
			return "MODEL.GETPOS", true, true
		case "SETROT", "ROT":
			return "MODEL.SETROT", true, true
		case "GETROT":
			return "MODEL.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "MODEL.SETSCALE", true, true
		case "GETSCALE":
			return "MODEL.GETSCALE", true, true
		case "SETCOLOR", "SETCOL", "COL":
			return "MODEL.SETCOLOR", true, true
		case "GETCOLOR", "GETCOL":
			return "MODEL.GETCOLOR", true, true
		case "SETALPHA", "ALPHA":
			return "MODEL.SETALPHA", true, true
		case "GETALPHA":
			return "MODEL.GETALPHA", true, true
		case "FREE":
			return "MODEL.FREE", true, true
		case "DRAW":
			return "MODEL.DRAW", true, true
		case "MOVE":
			return "MODEL.MOVE", true, true
		case "ROTATE":
			return "MODEL.ROTATE", true, true
		case "SETMATRIX", "MATRIX":
			return "MODEL.SETMATRIX", true, true
		case "SETSCALEUNIFORM":
			return "MODEL.SETSCALEUNIFORM", true, true
		case "SHOW", "HIDE":
			return "MODEL." + mn, true, true
		case "SETCULL", "CULL":
			return "MODEL.SETCULL", true, true
		case "SETWIREFRAME", "WIREFRAME":
			return "MODEL.SETWIREFRAME", true, true
		case "SETLIGHTING", "LIGHTING":
			return "MODEL.SETLIGHTING", true, true
		case "SETFOG", "FOG":
			return "MODEL.SETFOG", true, true
		case "SETBLEND", "BLEND":
			return "MODEL.SETBLEND", true, true
		case "SETDEPTH", "DEPTH":
			return "MODEL.SETDEPTH", true, true
		case "SETDIFFUSE", "DIFFUSE":
			return "MODEL.SETDIFFUSE", true, true
		case "SETSPECULAR", "SPECULAR":
			return "MODEL.SETSPECULAR", true, true
		case "SETSPECULARPOW", "SPECULARPOW":
			return "MODEL.SETSPECULARPOW", true, true
		case "SETEMISSIVE", "EMISSIVE":
			return "MODEL.SETEMISSIVE", true, true
		case "SETAMBIENTCOLOR", "AMBIENTCOLOR":
			return "MODEL.SETAMBIENTCOLOR", true, true
		case "SETMETAL", "METAL":
			return "MODEL.SETMETAL", true, true
		case "SETROUGH", "ROUGH":
			return "MODEL.SETROUGH", true, true
		case "SETGPUSKINNING", "GPUSKINNING":
			return "MODEL.SETGPUSKINNING", true, true
		case "SETTEXTURESTAGE", "TEXTURESTAGE":
			return "MODEL.SETTEXTURESTAGE", true, true
		case "SETSTAGEBLEND", "STAGEBLEND":
			return "MODEL.SETSTAGEBLEND", true, true
		case "SETSTAGESCROLL", "STAGESCROLL":
			return "MODEL.SETSTAGESCROLL", true, true
		case "SETSTAGESCALE", "STAGESCALE":
			return "MODEL.SETSTAGESCALE", true, true
		case "SETSTAGEROTATE", "STAGEROTATE":
			return "MODEL.SETSTAGEROTATE", true, true
		case "SCROLLTEXTURE":
			return "MODEL.SCROLLTEXTURE", true, true
		case "SCALETEXTURE":
			return "MODEL.SCALETEXTURE", true, true
		case "ROTATETEXTURE":
			return "MODEL.ROTATETEXTURE", true, true
		case "SETMATERIAL":
			return "MODEL.SETMATERIAL", true, true
		case "SETMATERIALTEXTURE", "MATERIALTEXTURE":
			return "MODEL.SETMATERIALTEXTURE", true, true
		case "SETMATERIALSHADER", "MATERIALSHADER":
			return "MODEL.SETMATERIALSHADER", true, true
		case "SETMODELMESHMATERIAL", "MODELMESHMATERIAL":
			return "MODEL.SETMODELMESHMATERIAL", true, true
		case "SETLODDISTANCES", "LODDISTANCES":
			return "MODEL.SETLODDISTANCES", true, true
		case "CLONE":
			return "MODEL.CLONE", true, true
		case "INSTANCE":
			return "MODEL.INSTANCE", true, true
		case "ATTACHTO", "ATTACH":
			return "MODEL.ATTACHTO", true, true
		case "DETACH":
			return "MODEL.DETACH", true, true
		case "EXISTS":
			return "MODEL.EXISTS", true, true
		}
	case heap.TagMaterial:
		switch mn {
		case "SETSHADER", "SHADER":
			return "MATERIAL.SETSHADER", true, true
		case "SETTEXTURE":
			return "MATERIAL.SETTEXTURE", true, true
		case "SETCOLOR", "SETCOL", "COL":
			return "MATERIAL.SETCOLOR", true, true
		case "SETFLOAT":
			return "MATERIAL.SETFLOAT", true, true
		case "SETEFFECT":
			return "MATERIAL.SETEFFECT", true, true
		case "SETEFFECTPARAM":
			return "MATERIAL.SETEFFECTPARAM", true, true
		case "FREE":
			return "MATERIAL.FREE", true, true
		}
	case heap.TagShader:
		switch mn {
		case "SETFLOAT", "FLOAT":
			return "SHADER.SETFLOAT", true, true
		case "SETVEC2", "VEC2":
			return "SHADER.SETVEC2", true, true
		case "SETVEC3", "VEC3":
			return "SHADER.SETVEC3", true, true
		case "SETVECTOR", "VECTOR":
			return "SHADER.SETVEC3", true, true
		case "SETVEC4", "VEC4":
			return "SHADER.SETVEC4", true, true
		case "SETINT", "INT":
			return "SHADER.SETINT", true, true
		case "SETTEXTURE":
			return "SHADER.SETTEXTURE", true, true
		case "GETLOC", "LOC":
			return "SHADER.GETLOC", true, true
		case "FREE":
			return "SHADER.FREE", true, true
		}
	case heap.TagComputeShader:
		switch mn {
		case "FREE":
			return "COMPUTESHADER.FREE", true, true
		case "SETBUFFER":
			return "COMPUTESHADER.SETBUFFER", true, true
		case "SETINT":
			return "COMPUTESHADER.SETINT", true, true
		case "SETFLOAT":
			return "COMPUTESHADER.SETFLOAT", true, true
		case "DISPATCH":
			return "COMPUTESHADER.DISPATCH", true, true
		}
	case heap.TagShaderBuffer:
		switch mn {
		case "FREE", "BUFFERFREE":
			return "COMPUTESHADER.BUFFERFREE", true, true
		}
	case heap.TagSky:
		switch mn {
		case "UPDATE":
			return "SKY.UPDATE", true, true
		case "DRAW":
			return "SKY.DRAW", true, true
		case "SETTIME", "TIME":
			return "SKY.SETTIME", true, true
		case "SETDAYLENGTH", "DAYLENGTH":
			return "SKY.SETDAYLENGTH", true, true
		case "GETTIMEHOURS":
			return "SKY.GETTIMEHOURS", true, true
		case "ISNIGHT":
			return "SKY.ISNIGHT", true, true
		case "FREE":
			return "SKY.FREE", true, true
		}
	case heap.TagCloud:
		switch mn {
		case "UPDATE":
			return "CLOUD.UPDATE", true, true
		case "DRAW":
			return "CLOUD.DRAW", true, true
		case "SETCOVERAGE", "COVERAGE":
			return "CLOUD.SETCOVERAGE", true, true
		case "GETCOVERAGE":
			return "CLOUD.GETCOVERAGE", true, true
		case "FREE":
			return "CLOUD.FREE", true, true
		}
	case heap.TagWeather:
		switch mn {
		case "UPDATE":
			return "WEATHER.UPDATE", true, true
		case "DRAW":
			return "WEATHER.DRAW", true, true
		case "SETTYPE", "TYPE":
			return "WEATHER.SETTYPE", true, true
		case "GETCOVERAGE":
			return "WEATHER.GETCOVERAGE", true, true
		case "GETTYPE":
			return "WEATHER.GETTYPE", true, true
		case "FREE":
			return "WEATHER.FREE", true, true
		}
	case heap.TagScatterSet:
		switch mn {
		case "APPLY":
			return "SCATTER.APPLY", true, true
		case "DRAWALL":
			return "SCATTER.DRAWALL", true, true
		case "FREE":
			return "SCATTER.FREE", true, true
		}
	case heap.TagProp:
		switch mn {
		case "PLACE":
			return "PROP.PLACE", true, true
		case "FREE":
			return "PROP.FREE", true, true
		}
	case heap.TagBrush:
		// Blitz-style entity brush API uses PascalCase registry keys (not BRUSH.*).
		switch mn {
		case "FREE", "FREEBRUSH":
			return "FreeBrush", true, true
		case "TEXTURE", "BRUSHTEXTURE":
			return "BrushTexture", true, true
		case "FX", "BRUSHFX":
			return "BrushFX", true, true
		case "SHININESS", "BRUSHSHININESS":
			return "BrushShininess", true, true
		case "SETCOLOR", "COLOR", "COL", "BRUSHCOLOR":
			return "BrushColor", true, true
		case "SETALPHA", "ALPHA", "BRUSHALPHA":
			return "BrushAlpha", true, true
		case "BLEND", "BRUSHBLEND":
			return "BrushBlend", true, true
		}
	case heap.TagMeshBuilder:
		switch mn {
		case "ADDVERTEX":
			return "ENTITY.ADDVERTEX", true, true
		case "ADDTRIANGLE":
			return "ENTITY.ADDTRIANGLE", true, true
		case "VERTEXX":
			return "ENTITY.VERTEXX", true, true
		case "VERTEXY":
			return "ENTITY.VERTEXY", true, true
		case "VERTEXZ":
			return "ENTITY.VERTEXZ", true, true
		case "X":
			return "ENTITY.VERTEXX", true, true
		case "Y":
			return "ENTITY.VERTEXY", true, true
		case "Z":
			return "ENTITY.VERTEXZ", true, true
		}
	case heap.TagTacticalGrid:
		switch mn {
		case "FREE":
			return "GRID.FREE", true, true
		case "SETCELL":
			return "GRID.SETCELL", true, true
		case "GETCELL":
			return "GRID.GETCELL", true, true
		case "WORLDTOCELL":
			return "GRID.WORLDTOCELL", true, true
		case "DRAW":
			return "GRID.DRAW", true, true
		case "SNAP":
			return "GRID.SNAP", true, true
		case "GETPATH":
			return "GRID.GETPATH", true, true
		case "FOLLOWTERRAIN":
			return "GRID.FOLLOWTERRAIN", true, true
		case "PLACEENTITY":
			return "GRID.PLACEENTITY", true, true
		case "RAYCAST":
			return "GRID.RAYCAST", true, true
		case "GETNEIGHBORS":
			return "GRID.GETNEIGHBORS", true, true
		}
	case heap.TagSpriteGroup:
		switch mn {
		case "ADD":
			return "SPRITEGROUP.ADD", true, true
		case "REMOVE":
			return "SPRITEGROUP.REMOVE", true, true
		case "CLEAR":
			return "SPRITEGROUP.CLEAR", true, true
		case "DRAW":
			return "SPRITEGROUP.DRAW", true, true
		case "FREE":
			return "SPRITEGROUP.FREE", true, true
		}
	case heap.TagSpriteLayer:
		switch mn {
		case "ADD":
			return "SPRITELAYER.ADD", true, true
		case "CLEAR":
			return "SPRITELAYER.CLEAR", true, true
		case "SETZ":
			return "SPRITELAYER.SETZ", true, true
		case "DRAW":
			return "SPRITELAYER.DRAW", true, true
		case "FREE":
			return "SPRITELAYER.FREE", true, true
		}
	case heap.TagSpriteBatch:
		switch mn {
		case "ADD":
			return "SPRITEBATCH.ADD", true, true
		case "CLEAR":
			return "SPRITEBATCH.CLEAR", true, true
		case "DRAW":
			return "SPRITEBATCH.DRAW", true, true
		case "FREE":
			return "SPRITEBATCH.FREE", true, true
		}
	case heap.TagSpriteUI:
		switch mn {
		case "DRAW":
			return "SPRITEUI.DRAW", true, true
		case "FREE":
			return "SPRITEUI.FREE", true, true
		}
	case heap.TagParticle2D:
		switch mn {
		case "EMIT":
			return "PARTICLE2D.EMIT", true, true
		case "UPDATE":
			return "PARTICLE2D.UPDATE", true, true
		case "DRAW":
			return "PARTICLE2D.DRAW", true, true
		case "FREE":
			return "PARTICLE2D.FREE", true, true
		}
	case heap.TagQuaternion:
		switch mn {
		case "FREE":
			return "QUAT.FREE", true, true
		case "NORMALIZE":
			return "QUAT.NORMALIZE", true, true
		case "INVERT":
			return "QUAT.INVERT", true, true
		case "TOMAT4":
			return "QUAT.TOMAT4", true, true
		case "TOEULER":
			return "QUAT.TOEULER", true, true
		case "MULTIPLY":
			return "QUAT.MULTIPLY", true, true
		case "SLERP":
			return "QUAT.SLERP", true, true
		case "TRANSFORM":
			return "QUAT.TRANSFORM", true, true
		}
	case heap.TagColor:
		switch mn {
		case "FREE":
			return "COLOR.FREE", true, true
		case "R":
			return "COLOR.R", true, true
		case "G":
			return "COLOR.G", true, true
		case "B":
			return "COLOR.B", true, true
		case "SETALPHA", "ALPHA":
			return "COLOR.A", true, true
		case "LERP":
			return "COLOR.LERP", true, true
		case "FADE":
			return "COLOR.FADE", true, true
		case "TOHSVX":
			return "COLOR.TOHSVX", true, true
		case "TOHSVY":
			return "COLOR.TOHSVY", true, true
		case "TOHSVZ":
			return "COLOR.TOHSVZ", true, true
		case "TOHSV":
			return "COLOR.TOHSV", true, true
		case "TOHEX":
			return "COLOR.TOHEX", true, true
		case "INVERT":
			return "COLOR.INVERT", true, true
		case "CONTRAST":
			return "COLOR.CONTRAST", true, true
		case "BRIGHTNESS":
			return "COLOR.BRIGHTNESS", true, true
		}
	case heap.TagVec2:
		switch mn {
		case "FREE":
			return "VEC2.FREE", true, true
		case "X":
			return "VEC2.X", true, true
		case "Y":
			return "VEC2.Y", true, true
		case "SET":
			return "VEC2.SET", true, true
		case "ADD":
			return "VEC2.ADD", true, true
		case "SUB":
			return "VEC2.SUB", true, true
		case "MUL":
			return "VEC2.MUL", true, true
		case "LEN", "LENGTH":
			return "VEC2.LENGTH", true, true
		case "NORMALIZE":
			return "VEC2.NORMALIZE", true, true
		case "LERP":
			return "VEC2.LERP", true, true
		case "DIST":
			return "VEC2.DIST", true, true
		case "DISTANCE":
			return "VEC2.DISTANCE", true, true
		case "ANGLE":
			return "VEC2.ANGLE", true, true
		case "ROTATE":
			return "VEC2.ROTATE", true, true
		case "TRANSFORMMAT4":
			return "VEC2.TRANSFORMMAT4", true, true
		}
	case heap.TagVec3:
		switch mn {
		case "FREE":
			return "VEC3.FREE", true, true
		case "X":
			return "VEC3.X", true, true
		case "Y":
			return "VEC3.Y", true, true
		case "Z":
			return "VEC3.Z", true, true
		case "SET":
			return "VEC3.SET", true, true
		case "ADD":
			return "VEC3.ADD", true, true
		case "SUB":
			return "VEC3.SUB", true, true
		case "MUL":
			return "VEC3.MUL", true, true
		case "DIV":
			return "VEC3.DIV", true, true
		case "DOT":
			return "VEC3.DOT", true, true
		case "CROSS":
			return "VEC3.CROSS", true, true
		case "LEN", "LENGTH":
			return "VEC3.LENGTH", true, true
		case "NORMALIZE":
			return "VEC3.NORMALIZE", true, true
		case "LERP":
			return "VEC3.LERP", true, true
		case "DIST":
			return "VEC3.DIST", true, true
		case "DISTANCE":
			return "VEC3.DISTANCE", true, true
		case "REFLECT":
			return "VEC3.REFLECT", true, true
		case "NEGATE":
			return "VEC3.NEGATE", true, true
		case "EQUALS":
			return "VEC3.EQUALS", true, true
		case "ANGLE":
			return "VEC3.ANGLE", true, true
		case "PROJECT":
			return "VEC3.PROJECT", true, true
		case "TRANSFORMMAT4":
			return "VEC3.TRANSFORMMAT4", true, true
		case "ROTATEBYQUAT":
			return "VEC3.ROTATEBYQUAT", true, true
		case "ORTHONORMALIZE":
			return "VEC3.ORTHONORMALIZE", true, true
		}
	case heap.TagAutomationList:
		switch mn {
		case "LISTEXPORT", "EXPORT":
			return "EVENT.LISTEXPORT", true, true
		case "SETACTIVELIST", "SETACTIVE":
			return "EVENT.SETACTIVELIST", true, true
		case "REPLAY":
			return "EVENT.REPLAY", true, true
		case "LISTCLEAR":
			return "EVENT.LISTCLEAR", true, true
		case "LISTCOUNT", "SETSIZE":
			return "EVENT.LISTCOUNT", true, true
		case "LISTFREE":
			return "EVENT.LISTFREE", true, true
		}
	case heap.TagBiome:
		switch mn {
		case "SETTEMP", "TEMP":
			return "BIOME.SETTEMP", true, true
		case "GETTEMP":
			return "BIOME.GETTEMP", true, true
		case "SETHUMIDITY", "HUMIDITY":
			return "BIOME.SETHUMIDITY", true, true
		case "GETHUMIDITY":
			return "BIOME.GETHUMIDITY", true, true
		case "FREE":
			return "BIOME.FREE", true, true
		}
	case heap.TagNoise:
		switch mn {
		case "FREE":
			return "NOISE.FREE", true, true
		case "SETTYPE", "TYPE":
			return "NOISE.SETTYPE", true, true
		case "SETSEED", "SEED":
			return "NOISE.SETSEED", true, true
		case "SETFREQUENCY", "FREQUENCY":
			return "NOISE.SETFREQUENCY", true, true
		case "SETOCTAVES", "OCTAVES":
			return "NOISE.SETOCTAVES", true, true
		case "SETLACUNARITY", "LACUNARITY":
			return "NOISE.SETLACUNARITY", true, true
		case "SETGAIN", "GAIN":
			return "NOISE.SETGAIN", true, true
		case "SETWEIGHTEDSTRENGTH":
			return "NOISE.SETWEIGHTEDSTRENGTH", true, true
		case "SETPINGPONGSTRENGTH":
			return "NOISE.SETPINGPONGSTRENGTH", true, true
		case "SETCELLULARTYPE":
			return "NOISE.SETCELLULARTYPE", true, true
		case "SETCELLULARDISTANCE":
			return "NOISE.SETCELLULARDISTANCE", true, true
		case "SETCELLULARJITTER":
			return "NOISE.SETCELLULARJITTER", true, true
		case "SETDOMAINWARPTYPE":
			return "NOISE.SETDOMAINWARPTYPE", true, true
		case "SETDOMAINWARPAMPLITUDE":
			return "NOISE.SETDOMAINWARPAMPLITUDE", true, true
		case "GET":
			return "NOISE.GET", true, true
		case "GET3D":
			return "NOISE.GET3D", true, true
		case "GETDOMAINWARPED":
			return "NOISE.GETDOMAINWARPED", true, true
		case "GETNORM":
			return "NOISE.GETNORM", true, true
		case "GETTILEABLE":
			return "NOISE.GETTILEABLE", true, true
		case "FILLARRAY":
			return "NOISE.FILLARRAY", true, true
		case "FILLARRAYNORM":
			return "NOISE.FILLARRAYNORM", true, true
		case "FILLIMAGE":
			return "NOISE.FILLIMAGE", true, true
		}
	case heap.TagTable:
		switch mn {
		case "FREE":
			return "TABLE.FREE", true, true
		case "ADDROW":
			return "TABLE.ADDROW", true, true
		case "ROWCOUNT", "ROWS", "SETSIZE":
			return "TABLE.ROWCOUNT", true, true
		case "COLCOUNT", "COLS", "COLUMNS":
			return "TABLE.COLCOUNT", true, true
		case "GET":
			return "TABLE.GET", true, true
		case "SET":
			return "TABLE.SET", true, true
		case "TOJSON":
			return "TABLE.TOJSON", true, true
		case "TOCSV":
			return "TABLE.TOCSV", true, true
		}
	case heap.TagPool:
		switch mn {
		case "SETFACTORY":
			return "POOL.SETFACTORY", true, true
		case "SETRESET":
			return "POOL.SETRESET", true, true
		case "PREWARM":
			return "POOL.PREWARM", true, true
		case "GET":
			return "POOL.GET", true, true
		case "RETURN":
			return "POOL.RETURN", true, true
		case "FREE":
			return "POOL.FREE", true, true
		}
	case heap.TagJSON:
		switch mn {
		case "FREE":
			return "JSON.FREE", true, true
		case "HAS":
			return "JSON.HAS", true, true
		case "TYPE":
			return "JSON.TYPE", true, true
		case "LEN", "LENGTH", "COUNT", "SETSIZE":
			return "JSON.LEN", true, true
		case "KEYS":
			return "JSON.KEYS", true, true
		case "GETSTRING":
			return "JSON.GETSTRING", true, true
		case "GETINT":
			return "JSON.GETINT", true, true
		case "GETFLOAT":
			return "JSON.GETFLOAT", true, true
		case "GETBOOL":
			return "JSON.GETBOOL", true, true
		case "GETARRAY":
			return "JSON.GETARRAY", true, true
		case "GETOBJECT":
			return "JSON.GETOBJECT", true, true
		case "SETSTRING":
			return "JSON.SETSTRING", true, true
		case "SETINT":
			return "JSON.SETINT", true, true
		case "SETFLOAT":
			return "JSON.SETFLOAT", true, true
		case "SETBOOL":
			return "JSON.SETBOOL", true, true
		case "SETNULL":
			return "JSON.SETNULL", true, true
		case "DELETE":
			return "JSON.DELETE", true, true
		case "CLEAR":
			return "JSON.CLEAR", true, true
		case "APPEND":
			return "JSON.APPEND", true, true
		case "TOSTRING":
			return "JSON.TOSTRING", true, true
		case "PRETTY":
			return "JSON.PRETTY", true, true
		case "MINIFY":
			return "JSON.MINIFY", true, true
		case "TOFILE":
			return "JSON.TOFILE", true, true
		case "SAVEFILE":
			return "JSON.SAVEFILE", true, true
		case "TOFILEPRETTY":
			return "JSON.TOFILEPRETTY", true, true
		case "TOCSV":
			return "JSON.TOCSV", true, true
		case "QUERY":
			return "JSON.QUERY", true, true
		}
	case heap.TagCSV:
		switch mn {
		case "FREE":
			return "CSV.FREE", true, true
		case "ROWCOUNT", "ROWS", "SETSIZE":
			return "CSV.ROWCOUNT", true, true
		case "COLCOUNT", "COLS", "COLUMNS":
			return "CSV.COLCOUNT", true, true
		case "GET":
			return "CSV.GET", true, true
		case "SET":
			return "CSV.SET", true, true
		case "TOSTRING":
			return "CSV.TOSTRING", true, true
		case "SAVE":
			return "CSV.SAVE", true, true
		case "TOJSON":
			return "CSV.TOJSON", true, true
		}
	case heap.TagDB:
		switch mn {
		case "CLOSE":
			return "DB.CLOSE", true, true
		case "ISOPEN":
			return "DB.ISOPEN", true, true
		case "EXEC":
			return "DB.EXEC", true, true
		case "QUERY":
			return "DB.QUERY", true, true
		case "QUERYJSON":
			return "DB.QUERYJSON", true, true
		case "PREPARE":
			return "DB.PREPARE", true, true
		case "BEGIN":
			return "DB.BEGIN", true, true
		case "LASTINSERTID":
			return "DB.LASTINSERTID", true, true
		case "CHANGES":
			return "DB.CHANGES", true, true
		}
	case heap.TagDBRows:
		switch mn {
		case "NEXT":
			return "ROWS.NEXT", true, true
		case "CLOSE":
			return "ROWS.CLOSE", true, true
		case "GETSTRING":
			return "ROWS.GETSTRING", true, true
		case "GETINT":
			return "ROWS.GETINT", true, true
		case "GETFLOAT":
			return "ROWS.GETFLOAT", true, true
		}
	case heap.TagDBStmt:
		switch mn {
		case "STMTCLOSE", "CLOSE":
			return "DB.STMTCLOSE", true, true
		case "STMTEXEC", "EXEC":
			return "DB.STMTEXEC", true, true
		}
	case heap.TagDBTx:
		switch mn {
		case "COMMIT":
			return "DB.COMMIT", true, true
		case "ROLLBACK":
			return "DB.ROLLBACK", true, true
		}
	case heap.TagRng:
		switch mn {
		case "NEXT":
			return "RAND.NEXT", true, true
		case "NEXTF":
			return "RAND.NEXTF", true, true
		case "FREE":
			return "RAND.FREE", true, true
		}
	case heap.TagMem:
		switch mn {
		case "FREE":
			return "MEM.FREE", true, true
		case "SIZE", "SETSIZE":
			return "MEM.SIZE", true, true
		case "CLEAR":
			return "MEM.CLEAR", true, true
		case "GETBYTE":
			return "MEM.GETBYTE", true, true
		case "GETWORD":
			return "MEM.GETWORD", true, true
		case "GETDWORD":
			return "MEM.GETDWORD", true, true
		case "GETFLOAT":
			return "MEM.GETFLOAT", true, true
		case "GETDOUBLE":
			return "MEM.GETDOUBLE", true, true
		case "GETSTRING":
			return "MEM.GETSTRING", true, true
		case "SETBYTE":
			return "MEM.SETBYTE", true, true
		case "SETWORD":
			return "MEM.SETWORD", true, true
		case "SETDWORD":
			return "MEM.SETDWORD", true, true
		case "SETFLOAT":
			return "MEM.SETFLOAT", true, true
		case "SETDOUBLE":
			return "MEM.SETDOUBLE", true, true
		case "SETSTRING":
			return "MEM.SETSTRING", true, true
		case "RESIZE":
			return "MEM.RESIZE", true, true
		}
	case heap.TagLobby:
		switch mn {
		case "FREE":
			return "LOBBY.FREE", true, true
		case "SETPROPERTY":
			return "LOBBY.SETPROPERTY", true, true
		case "SETHOST":
			return "LOBBY.SETHOST", true, true
		case "START":
			return "LOBBY.START", true, true
		// LOBBY.FIND is global (key$, value$) — not receiver-first; use LOBBY.FIND(...)
		case "GETNAME":
			return "LOBBY.GETNAME", true, true
		case "JOIN":
			return "LOBBY.JOIN", true, true
		}
	case heap.TagNetPacket:
		switch mn {
		case "DATA":
			return "PACKET.DATA", true, true
		case "FREE":
			return "PACKET.FREE", true, true
		}
	case heap.TagHost:
		switch mn {
		case "UPDATE", "PUMP":
			return "NET.UPDATE", true, true
		case "RECEIVE", "POP":
			return "NET.RECEIVE", true, true
		case "CLOSE":
			return "NET.CLOSE", true, true
		case "BROADCAST":
			return "NET.BROADCAST", true, true
		case "PEERCOUNT", "PEERS":
			return "NET.PEERCOUNT", true, true
		case "SETBANDWIDTH", "BANDWIDTH":
			return "NET.SETBANDWIDTH", true, true
		case "SERVICE":
			return "NET.SERVICE", true, true
		case "FLUSH":
			return "NET.FLUSH", true, true
		}
	case heap.TagEvent:
		switch mn {
		case "TYPE":
			return "EVENT.TYPE", true, true
		case "PEER":
			return "EVENT.PEER", true, true
		case "DATA":
			return "EVENT.DATA", true, true
		case "FREE":
			return "EVENT.FREE", true, true
		case "CHANNEL":
			return "EVENT.CHANNEL", true, true
		}
	case heap.TagPlayer2D:
		switch mn {
		case "FREE":
			return "PLAYER2D.FREE", true, true
		case "MOVE":
			return "PLAYER2D.MOVE", true, true
		case "CLAMP":
			return "PLAYER2D.CLAMP", true, true
		case "KEEPINBOUNDS", "BOUNDS":
			return "PLAYER2D.KEEPINBOUNDS", true, true
		case "GETX", "X":
			return "PLAYER2D.GETX", true, true
		case "GETZ", "Z":
			return "PLAYER2D.GETZ", true, true
		case "GETPOS":
			return "PLAYER2D.GETPOS", true, true
		case "SETPOS", "POS":
			return "PLAYER2D.SETPOS", true, true
		}
	case heap.TagGameTimer:
		switch mn {
		case "RESET":
			return "TIMER.RESET", true, true
		case "FINISHED", "DONE":
			return "TIMER.FINISHED", true, true
		case "FREE":
			return "TIMER.FREE", true, true
		case "REMAINING", "LEFT":
			return "TIMER.REMAINING", true, true
		}
	case heap.TagGameTimerSim:
		switch mn {
		case "START":
			return "TIMER.START", true, true
		case "STOP":
			return "TIMER.STOP", true, true
		case "REWIND":
			return "TIMER.REWIND", true, true
		case "SETLOOP", "LOOP":
			return "TIMER.SETLOOP", true, true
		case "GETLOOP":
			return "TIMER.GETLOOP", true, true
		case "UPDATE":
			return "TIMER.UPDATE", true, true
		case "DONE":
			return "TIMER.DONE", true, true
		case "FRACTION", "FRAC":
			return "TIMER.FRACTION", true, true
		case "FREE":
			return "TIMER.FREE", true, true
		case "REMAINING", "LEFT":
			return "TIMER.REMAINING", true, true
		}
	case heap.TagGameStopwatch:
		switch mn {
		case "RESET":
			return "STOPWATCH.RESET", true, true
		case "ELAPSED", "TIME":
			return "STOPWATCH.ELAPSED", true, true
		case "FREE":
			return "STOPWATCH.FREE", true, true
		}
	case heap.TagDecal:
		switch mn {
		case "SETPOS", "POS":
			return "DECAL.SETPOS", true, true
		case "GETPOS":
			return "DECAL.GETPOS", true, true
		case "SETROT", "ROT":
			return "DECAL.SETROT", true, true
		case "GETROT":
			return "DECAL.GETROT", true, true
		case "SETSIZE", "SIZE":
			return "DECAL.SETSIZE", true, true
		case "GETSIZE":
			return "DECAL.GETSIZE", true, true
		case "SETCOLOR", "COLOR":
			return "DECAL.SETCOLOR", true, true
		case "GETCOLOR":
			return "DECAL.GETCOLOR", true, true
		case "SETALPHA", "ALPHA":
			return "DECAL.SETALPHA", true, true
		case "GETALPHA":
			return "DECAL.GETALPHA", true, true
		case "SETLIFETIME", "LIFETIME":
			return "DECAL.SETLIFETIME", true, true
		case "GETLIFETIME":
			return "DECAL.GETLIFETIME", true, true
		case "DRAW":
			return "DECAL.DRAW", true, true
		case "FREE":
			return "DECAL.FREE", true, true
		}
	case heap.TagWater:
		switch mn {
		case "SETPOS", "POS":
			return "WATER.SETPOS", true, true
		case "GETPOS":
			return "WATER.GETPOS", true, true
		case "SETROT", "ROT":
			return "WATER.SETROT", true, true
		case "GETROT":
			return "WATER.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "WATER.SETSCALE", true, true
		case "GETSCALE":
			return "WATER.GETSCALE", true, true
		case "DRAW":
			return "WATER.DRAW", true, true
		case "FREE":
			return "WATER.FREE", true, true
		case "SETWAVE", "WAVE":
			return "WATER.SETWAVE", true, true
		case "SETWAVEHEIGHT", "WAVEHEIGHT":
			return "WATER.SETWAVEHEIGHT", true, true
		case "GETWAVEHEIGHT":
			return "WATER.GETWAVEHEIGHT", true, true
		case "GETWAVESPEED":
			return "WATER.GETWAVESPEED", true, true
		case "SETCOLOR", "SETCOL", "COL":
			return "WATER.SETCOLOR", true, true
		case "GETCOLOR", "GETCOL":
			return "WATER.GETCOLOR", true, true
		case "SETSHALLOWCOLOR", "SHALLOWCOLOR":
			return "WATER.SETSHALLOWCOLOR", true, true
		case "SETDEEPCOLOR", "DEEPCOLOR":
			return "WATER.SETDEEPCOLOR", true, true
		case "GETSHALLOWCOLOR":
			return "WATER.GETSHALLOWCOLOR", true, true
		case "GETDEEPCOLOR":
			return "WATER.GETDEEPCOLOR", true, true
		case "GETWAVEY":
			return "WATER.GETWAVEY", true, true
		case "GETDEPTH":
			return "WATER.GETDEPTH", true, true
		case "ISUNDER":
			return "WATER.ISUNDER", true, true
		case "AUTOPHYSICS":
			return "WATER.AUTOPHYSICS", true, true
		case "UPDATE":
			return "WATER.UPDATE", true, true
		}
	case heap.TagParticle:
		switch mn {
		case "SETPOS", "POS":
			return "PARTICLE.SETPOS", true, true
		case "GETPOS":
			return "PARTICLE.GETPOS", true, true
		case "SETCOLOR", "SETCOL", "COL":
			return "PARTICLE.SETCOLOR", true, true
		case "GETCOLOR", "GETCOL":
			return "PARTICLE.GETCOLOR", true, true
		case "SETALPHA", "ALPHA":
			return "PARTICLE.SETALPHA", true, true
		case "GETALPHA":
			return "PARTICLE.GETALPHA", true, true
		case "GETVELOCITY":
			return "PARTICLE.GETVELOCITY", true, true
		case "GETSIZE":
			return "PARTICLE.GETSIZE", true, true
		case "SETVELOCITY", "VEL":
			return "PARTICLE.SETVELOCITY", true, true
		// normalize maps scale → SETSCALE; there is no PARTICLE.SETSCALE (emitter size uses SETSIZE).
		case "SETSCALE":
			return "PARTICLE.SETSIZE", true, true
		case "SETTEXTURE", "SETEMITRATE", "SETRATE", "SETLIFETIME", "SETDIRECTION", "SETSPREAD", "SETSPEED", "SETSTARTSIZE", "SETENDSIZE", "SETSIZE", "SETGRAVITY", "SETBURST", "SETBILLBOARD", "PLAY", "STOP", "UPDATE", "DRAW", "ISALIVE", "COUNT", "FREE":
			return "PARTICLE." + mn, true, true
		}
	case heap.TagTerrain:
		switch mn {
		case "SETPOS", "POS":
			return "TERRAIN.SETPOS", true, true
		case "GETPOS":
			return "TERRAIN.GETPOS", true, true
		case "SETROT", "ROT":
			return "TERRAIN.SETROT", true, true
		case "GETROT":
			return "TERRAIN.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "TERRAIN.SETSCALE", true, true
		case "GETSCALE":
			return "TERRAIN.GETSCALE", true, true
		case "FREE":
			return "TERRAIN.FREE", true, true
		case "APPLYMAP":
			return "TERRAIN.APPLYMAP", true, true
		case "APPLYTILES":
			return "TERRAIN.APPLYTILES", true, true
		case "DRAW":
			return "TERRAIN.DRAW", true, true
		case "SETCHUNKSIZE":
			return "TERRAIN.SETCHUNKSIZE", true, true
		case "FILLPERLIN":
			return "TERRAIN.FILLPERLIN", true, true
		case "FILLFLAT":
			return "TERRAIN.FILLFLAT", true, true
		case "RAISE":
			return "TERRAIN.RAISE", true, true
		case "LOWER":
			return "TERRAIN.LOWER", true, true
		case "SETMESHBUILDBUDGET":
			return "TERRAIN.SETMESHBUILDBUDGET", true, true
		case "SETASYNCMESHBUILD", "ASYNC", "ASYNCMESH":
			return "TERRAIN.SETASYNCMESHBUILD", true, true
		case "LOAD":
			return "TERRAIN.LOAD", true, true
		case "SETDETAIL":
			return "TERRAIN.SETDETAIL", true, true
		case "GETDETAIL":
			return "TERRAIN.GETDETAIL", true, true
		case "GETHEIGHT":
			return "TERRAIN.GETHEIGHT", true, true
		case "GETSLOPE":
			return "TERRAIN.GETSLOPE", true, true
		case "GETNORMAL":
			return "TERRAIN.GETNORMAL", true, true
		case "GETSPLAT":
			return "TERRAIN.GETSPLAT", true, true
		case "RAYCAST":
			return "TERRAIN.RAYCAST", true, true
		case "PLACE":
			return "TERRAIN.PLACE", true, true
		case "SNAPPY":
			return "TERRAIN.SNAPY", true, true
		}
	case heap.TagInstancedModel:
		switch mn {
		case "SETPOS", "POS":
			return "INSTANCE.SETPOS", true, true
		case "SETINSTANCEPOS":
			return "INSTANCE.SETINSTANCEPOS", true, true
		case "GETPOS":
			return "INSTANCE.GETPOS", true, true
		case "SETROT", "ROT":
			return "INSTANCE.SETROT", true, true
		case "GETROT":
			return "INSTANCE.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "INSTANCE.SETSCALE", true, true
		case "SETINSTANCESCALE":
			return "INSTANCE.SETINSTANCESCALE", true, true
		case "GETSCALE":
			return "INSTANCE.GETSCALE", true, true
		case "SETCOLOR", "SETCOL", "COL":
			return "INSTANCE.SETCOLOR", true, true
		case "GETCOLOR", "GETCOL":
			return "INSTANCE.GETCOLOR", true, true
		case "GETALPHA":
			return "INSTANCE.GETALPHA", true, true
		case "SETMATRIX":
			return "INSTANCE.SETMATRIX", true, true
		case "UPDATEINSTANCES", "UPDATEBUFFER":
			return "INSTANCE.UPDATEINSTANCES", true, true
		case "SETCULLDISTANCE":
			return "INSTANCE.SETCULLDISTANCE", true, true
		case "DRAWLOD":
			return "INSTANCE.DRAWLOD", true, true
		case "DRAW":
			return "INSTANCE.DRAW", true, true
		case "FREE":
			return "INSTANCE.FREE", true, true
		case "COUNT":
			return "INSTANCE.COUNT", true, true
		}
	case heap.TagMatrix:
		switch mn {
		case "FREE":
			return "TRANSFORM.FREE", true, true
		// normalizeHandleMethod maps SETROTATION / ROT / … → SETROT
		case "SETROT":
			return "TRANSFORM.SETROTATION", true, true
		case "INVERSE":
			return "TRANSFORM.INVERSE", true, true
		case "TRANSPOSE":
			return "TRANSFORM.TRANSPOSE", true, true
		case "MULTIPLY":
			return "TRANSFORM.MULTIPLY", true, true
		case "GETELEMENT":
			return "TRANSFORM.GETELEMENT", true, true
		case "APPLYX":
			return "TRANSFORM.APPLYX", true, true
		case "APPLYY":
			return "TRANSFORM.APPLYY", true, true
		case "APPLYZ":
			return "TRANSFORM.APPLYZ", true, true
		// MAT4.TRANSFORM* aliases (same implementation as TRANSFORM.APPLY*).
		case "TRANSFORMX":
			return "TRANSFORM.APPLYX", true, true
		case "TRANSFORMY":
			return "TRANSFORM.APPLYY", true, true
		case "TRANSFORMZ":
			return "TRANSFORM.APPLYZ", true, true
		// Rotation extraction: matrix handle is first argument to QUAT.FROMMAT4.
		case "TOQUAT", "FROMMAT4":
			return "QUAT.FROMMAT4", true, true
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
	case heap.TagTexture, heap.TagImage, heap.TagImageSequence:
		pre := "TEXTURE"
		if tag == heap.TagImage || tag == heap.TagImageSequence {
			pre = "IMAGE"
		}
		switch mn {
		case "WIDTH", "GETWIDTH":
			return pre + ".GETWIDTH", true, true
		case "HEIGHT", "GETHEIGHT":
			return pre + ".GETHEIGHT", true, true
		case "SIZE", "GETSIZE":
			return pre + ".GETSIZE", true, true
		case "FREE":
			return pre + ".FREE", true, true
		case "CROP", "RESIZE", "FLIPH", "FLIPV", "ROTATE", "ROTATECW", "ROTATECCW", "COLORTINT", "COLORINVERT", "COLORGRAYSCALE", "COLORCONTRAST", "COLORBRIGHTNESS", "COLORREPLACE", "CLEAR":
			if tag == heap.TagImage {
				return "IMAGE." + mn, true, true
			}
		}
	case heap.TagFile:
		switch mn {
		case "SETPOS", "POS", "SEEK":
			return "FILE.SEEK", true, true
		case "TELL":
			return "FILE.TELL", true, true
		case "GETPOS":
			return "FILE.GETPOS", true, true
		case "GETSIZE", "SIZE", "SETSIZE":
			return "FILE.GETSIZE", true, true
		case "EOF", "GETEOF":
			return "FILE.GETEOF", true, true
		case "READLINE", "READ":
			return "FILE.READLINE", true, true
		case "READSTRING":
			return "FILE.READSTRING", true, true
		case "WRITE", "WRITELN":
			return "FILE.WRITE", true, true
		case "CLOSE":
			return "FILE.CLOSE", true, true
		}
	case heap.TagArray:
		switch mn {
		case "LEN", "COUNT", "SIZE", "GETLEN", "SETSIZE":
			return "ARRAY.LEN", true, true
		case "FILL", "SETALL":
			return "ARRAY.FILL", true, true
		case "COPY", "SET":
			return "ARRAY.COPY", true, true
		case "SORT":
			return "ARRAY.SORT", true, true
		case "REVERSE":
			return "ARRAY.REVERSE", true, true
		case "PUSH":
			return "ARRAY.PUSH", true, true
		case "POP":
			return "ARRAY.POP", true, true
		case "SHIFT":
			return "ARRAY.SHIFT", true, true
		case "UNSHIFT":
			return "ARRAY.UNSHIFT", true, true
		case "SPLICE":
			return "ARRAY.SPLICE", true, true
		case "SLICE":
			return "ARRAY.SLICE", true, true
		case "JOIN":
			return "ARRAY.JOINS", true, true
		case "FIND":
			return "ARRAY.FIND", true, true
		case "CONTAINS":
			return "ARRAY.CONTAINS", true, true
		case "FREE":
			return "ARRAY.FREE", true, true
		}
	case heap.TagSound, heap.TagMusic:
		switch mn {
		case "PLAY":
			return "AUDIO.PLAY", true, true
		case "STOP":
			return "AUDIO.STOP", true, true
		case "PAUSE":
			return "AUDIO.PAUSE", true, true
		case "RESUME":
			return "AUDIO.RESUME", true, true
		case "UPDATE":
			if tag == heap.TagMusic {
				return "AUDIO.UPDATEMUSIC", true, true
			}
			return "", false, false
		case "SETVOLUME", "VOLUME":
			if tag == heap.TagSound {
				return "AUDIO.SETSOUNDVOLUME", true, true
			}
			return "AUDIO.SETMUSICVOLUME", true, true
		case "SETPITCH", "PITCH":
			if tag == heap.TagSound {
				return "AUDIO.SETSOUNDPITCH", true, true
			}
			return "AUDIO.SETMUSICPITCH", true, true
		case "SETPAN", "PAN":
			return "AUDIO.SETSOUNDPAN", true, true
		case "SEEK":
			return "AUDIO.SEEKMUSIC", true, true
		case "GETVOLUME":
			if tag == heap.TagSound {
				return "AUDIO.GETSOUNDVOLUME", true, true
			}
			return "AUDIO.GETMUSICVOLUME", true, true
		case "GETPITCH":
			if tag == heap.TagSound {
				return "AUDIO.GETSOUNDPITCH", true, true
			}
			return "AUDIO.GETMUSICPITCH", true, true
		case "GETPAN":
			if tag == heap.TagSound {
				return "AUDIO.GETSOUNDPAN", true, true
			}
			return "", false, false
		case "LENGTH":
			if tag == heap.TagMusic {
				return "AUDIO.GETMUSICLENGTH", true, true
			}
			return "", false, false
		case "TIME":
			if tag == heap.TagMusic {
				return "AUDIO.GETMUSICTIME", true, true
			}
			return "", false, false
		case "FREE":
			if tag == heap.TagSound {
				return "SOUND.FREE", true, true
			}
			return "MUSIC.FREE", true, true
		}
	case heap.TagAudioStream:
		switch mn {
		case "PLAY":
			return "AUDIOSTREAM.PLAY", true, true
		case "STOP":
			return "AUDIOSTREAM.STOP", true, true
		case "PAUSE":
			return "AUDIOSTREAM.PAUSE", true, true
		case "RESUME":
			return "AUDIOSTREAM.RESUME", true, true
		case "SETVOLUME", "VOLUME":
			return "AUDIOSTREAM.SETVOLUME", true, true
		case "GETVOLUME":
			return "AUDIOSTREAM.GETVOLUME", true, true
		case "SETPITCH", "PITCH":
			return "AUDIOSTREAM.SETPITCH", true, true
		case "GETPITCH":
			return "AUDIOSTREAM.GETPITCH", true, true
		case "SETPAN", "PAN":
			return "AUDIOSTREAM.SETPAN", true, true
		case "GETPAN":
			return "AUDIOSTREAM.GETPAN", true, true
		case "UPDATE":
			return "AUDIOSTREAM.UPDATE", true, true
		case "FREE":
			return "AUDIOSTREAM.FREE", true, true
		}
	case heap.TagWave:
		switch mn {
		case "COPY":
			return "WAVE.COPY", true, true
		case "CROP":
			return "WAVE.CROP", true, true
		case "FORMAT":
			return "WAVE.FORMAT", true, true
		case "EXPORT":
			return "WAVE.EXPORT", true, true
		case "FREE":
			return "WAVE.FREE", true, true
		}
	case heap.TagFont:
		switch mn {
		// Text metrics use DRAW.TEXTWIDTH / DRAW.TEXTFONTWIDTH (font is not a standalone measure API today).
		case "SETDEFAULT", "DEFAULT":
			return "FONT.SETDEFAULT", true, true
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

// handleCallDispatch resolves handle method calls. When argCount==0, universal pose setters
// (pos/rot/scale) dispatch to GET* builtins so e.g. handle.pos() reads position.
func handleCallDispatch(tag uint16, method string, argCount int) (registryKey string, prependReceiver bool, ok bool) {
	if argCount > 0 {
		return handleCallBuiltin(tag, method)
	}
	mn := normalizeHandleMethod(strings.ToUpper(strings.TrimSpace(method)))
	switch tag {
	case heap.TagCamera:
		switch mn {
		case "SETPOS", "POS", "POSITION":
			return "CAMERA.GETPOS", true, true
		case "SETROT", "ROT":
			return "CAMERA.GETROT", true, true
		case "SETFOV", "FOV":
			return "CAMERA.GETFOV", true, true
		case "SETTARGET", "TARGET":
			return "CAMERA.GETTARGET", true, true
		case "SETUP", "UP":
			return "CAMERA.GETUP", true, true
		case "PROJECTION", "SETPROJECTION":
			return "CAMERA.GETPROJECTION", true, true
		}
	case heap.TagModel, heap.TagLODModel:
		switch mn {
		case "SETPOS", "POS":
			return "MODEL.GETPOS", true, true
		case "SETROT", "ROT":
			return "MODEL.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "MODEL.GETSCALE", true, true
		case "SETCOLOR", "COL":
			return "MODEL.GETCOLOR", true, true
		case "SETALPHA", "ALPHA":
			return "MODEL.GETALPHA", true, true
		}
	case heap.TagLight:
		switch mn {
		case "SETPOS", "POS", "POSITION":
			return "LIGHT.GETPOS", true, true
		case "SETROT", "ROT", "SETDIR", "DIR":
			return "LIGHT.GETDIR", true, true
		case "SETCOLOR", "COL", "COLOR":
			return "LIGHT.GETCOLOR", true, true
		case "SETINTENSITY", "SETENERGY", "INTENSITY", "ENERGY", "GETENERGY": // Pure getters
			return "LIGHT.GETINTENSITY", true, true
		case "RANGE", "SETRANGE":
			return "LIGHT.GETRANGE", true, true
		case "SHADOW", "SETSHADOW":
			return "LIGHT.GETSHADOW", true, true
		case "GETPOS", "GETDIR", "GETCOLOR", "GETINTENSITY", "GETSHADOW", "GETRANGE", "GETINNERCONE", "GETOUTERCONE": // Pure getters
			return "LIGHT." + mn, true, true
		case "ENABLED", "ENABLE", "ISENABLED":
			return "LIGHT.ISENABLED", true, true
		}
	case heap.TagParticle:
		switch mn {
		case "SETPOS", "SETPOSITION", "POS", "POSITION":
			return "PARTICLE.GETPOS", true, true
		case "SETCOLOR", "COLOR", "COL":
			return "PARTICLE.GETCOLOR", true, true
		case "SETALPHA", "ALPHA", "A":
			return "PARTICLE.GETALPHA", true, true
		case "SETVELOCITY", "VEL":
			return "PARTICLE.GETVELOCITY", true, true
		// normalize maps size/scale → SETSIZE / SETSCALE; writes use SETSIZE (see handleCallBuiltin).
		case "SETSIZE", "SIZE", "SETSCALE", "SCALE":
			return "PARTICLE.GETSIZE", true, true
		case "GETPOS", "GETCOLOR", "GETALPHA", "GETVELOCITY", "GETSIZE":
			return "PARTICLE." + mn, true, true
		}

	case heap.TagInstancedModel:
		switch mn {
		case "SETPOS", "POS":
			return "INSTANCE.GETPOS", true, true
		case "SETROT", "ROT":
			return "INSTANCE.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "INSTANCE.GETSCALE", true, true
		case "SETCOLOR", "COL":
			return "INSTANCE.GETCOLOR", true, true
		case "SETALPHA", "ALPHA":
			return "INSTANCE.GETALPHA", true, true
		case "GETPOS", "GETROT", "GETSCALE", "GETCOLOR", "GETALPHA":
			return "INSTANCE." + mn, true, true
		}
	case heap.TagPhysicsBody:
		switch mn {
		case "SETPOS", "POS":
			return "BODY3D.GETPOS", true, true
		case "SETROT", "ROT":
			return "BODY3D.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "BODY3D.GETSCALE", true, true
		case "SETVELOCITY", "SETVEL", "VEL", "VELOCITY", "GETVELOCITY", "GETVEL":
			return "BODY3D.GETVELOCITY", true, true
		case "SETANGULARVELOCITY", "SETANGULARVEL", "ANGULARVELOCITY", "ANGULARVEL", "ANGVEL", "GETANGULARVEL":
			return "BODY3D.GETANGULARVEL", true, true
		case "SETFRICTION", "FRICTION":
			return "BODY3D.GETFRICTION", true, true
		case "SETRESTITUTION", "RESTITUTION", "SETBOUNCE", "BOUNCE":
			return "BODY3D.GETRESTITUTION", true, true
		case "SETGRAVITYFACTOR", "GRAVITYFACTOR":
			return "BODY3D.GETGRAVITYFACTOR", true, true
		case "SETDAMPING", "DAMPING":
			return "BODY3D.GETDAMPING", true, true
		case "SETCCD", "CCD":
			return "BODY3D.GETCCD", true, true
		case "SETMASS", "MASS", "GETMASS":
			return "BODY3D.GETMASS", true, true
		case "X":
			return "BODY3D.X", true, true
		case "Y":
			return "BODY3D.Y", true, true
		case "Z":
			return "BODY3D.Z", true, true
		case "FREE":
			return "BODY3D.FREE", true, true
		}
	case heap.TagBody2D:
		switch mn {
		case "SETPOS", "POS":
			return "BODY2D.GETPOS", true, true
		case "SETROT", "ROT":
			return "BODY2D.GETROT", true, true
		case "SETVELOCITY", "SETVEL", "VEL", "LINEARVEL", "VELOCITY":
			return "BODY2D.GETLINEARVELOCITY", true, true
		case "SETANGULARVELOCITY", "SETANGULARVEL", "ANGULARVELOCITY", "ANGULARVEL", "ANGVEL":
			return "BODY2D.GETANGULARVELOCITY", true, true
		case "SETMASS", "MASS", "GETMASS":
			return "BODY2D.GETMASS", true, true
		case "FRICTION", "SETFRICTION":
			return "BODY2D.GETFRICTION", true, true
		case "RESTITUTION", "SETRESTITUTION", "BOUNCE":
			return "BODY2D.GETRESTITUTION", true, true
		}
	case heap.TagCharController:
		switch mn {
		case "SETPOS", "POS":
			return "CHARACTERREF.GETPOSITION", true, true
		case "SETROT", "ROT":
			return "CHARACTERREF.GETROT", true, true
		case "SETMAXSLOPE", "MAXSLOPE", "SLOPE", "SETSLOPE":
			return "CHARACTERREF.GETMAXSLOPE", true, true
		case "SETSTEPHEIGHT", "STEPHEIGHT", "STEP", "SETSTEP":
			return "CHARACTERREF.GETSTEPHEIGHT", true, true
		case "SETSNAPDISTANCE", "SNAPDISTANCE", "SNAP", "SETSTICKDOWN":
			return "CHARACTERREF.GETSNAPDISTANCE", true, true
		case "SETGRAVITY", "GRAVITY", "SETGRAVITYSCALE", "GRAVITYSCALE":
			return "CHARACTERREF.GETGRAVITY", true, true
		case "SETFRICTION", "FRICTION":
			return "CHARACTERREF.GETFRICTION", true, true
		case "SETBOUNCE", "SETBOUNCINESS", "BOUNCE", "BOUNCINESS", "GETBOUNCE", "GETBOUNCINESS":
			return "CHARACTERREF.GETBOUNCINESS", true, true
		case "SETPADDING", "PADDING":
			return "CHARACTERREF.GETPADDING", true, true
		case "ISGROUNDED", "GROUNDED":
			return "CHARACTERREF.ISGROUNDED", true, true
		case "GETSLOPEANGLE", "SLOPEANGLE":
			return "CHARACTERREF.GETSLOPEANGLE", true, true
		case "GETSPEED", "SETSPEED", "SPEED":
			return "CHARACTERREF.GETSPEED", true, true
		case "SETVELOCITY", "SETVEL", "GETVELOCITY", "VELOCITY", "VEL":
			return "CHARACTERREF.GETVELOCITY", true, true
		case "ISMOVING", "MOVING":
			return "CHARACTERREF.ISMOVING", true, true
		case "UPDATE":
			return "CHARACTERREF.UPDATE", true, true
		case "X":
			return "CHARCONTROLLER.X", true, true
		case "Y":
			return "CHARCONTROLLER.Y", true, true
		case "Z":
			return "CHARCONTROLLER.Z", true, true
		case "FREE":
			return "CHARACTERREF.FREE", true, true
		case "SETJUMPBUFFER", "JUMPBUFFER", "GETJUMPBUFFER":
			return "CHARACTERREF.GETJUMPBUFFER", true, true
		case "SETAIRCONTROL", "AIRCONTROL", "GETAIRCONTROL":
			return "CHARACTERREF.GETAIRCONTROL", true, true
		case "SETGROUNDCONTROL", "GROUNDCONTROL", "GETGROUNDCONTROL":
			return "CHARACTERREF.GETGROUNDCONTROL", true, true
		}
	case heap.TagTween:
		switch mn {
		case "ISPLAYING", "PLAYING":
			return "TWEEN.ISPLAYING", true, true
		case "ISFINISHED", "FINISHED":
			return "TWEEN.ISFINISHED", true, true
		case "PROGRESS":
			return "TWEEN.PROGRESS", true, true
		case "SETLOOP", "LOOP", "GETLOOP":
			return "TWEEN.GETLOOP", true, true
		case "YOYO", "GETYOYO":
			return "TWEEN.GETYOYO", true, true
		}
	case heap.TagRay:
		switch mn {
		// normalize maps "pos"/"position" → SETPOS; must read, not RAY.SETPOS.
		case "SETPOS", "POS", "POSITION":
			return "RAY.GETPOS", true, true
		case "SETDIR", "DIR", "DIRECTION":
			return "RAY.GETDIR", true, true
		}
	case heap.TagBBox:
		switch mn {
		case "SETMIN", "MIN":
			return "BBOX.GETMIN", true, true
		case "SETMAX", "MAX":
			return "BBOX.GETMAX", true, true
		}
	case heap.TagBSphere:
		switch mn {
		case "SETPOS", "POS", "POSITION":
			return "BSPHERE.GETPOS", true, true
		case "SETRADIUS", "RADIUS":
			return "BSPHERE.GETRADIUS", true, true
		}
	case heap.TagEntityRef:
		switch mn {
		case "SETPOS", "POS":
			return "ENTITY.GETPOS", true, true
		case "SETROT", "ROT":
			return "ENTITY.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "ENTITY.GETSCALE", true, true
		case "SETCOLOR", "COLOR", "COL":
			return "ENTITY.GETCOLOR", true, true
		case "SETALPHA", "ALPHA":
			return "ENTITY.GETALPHA", true, true
		case "GETPOS", "GETROT", "GETSCALE", "GETCOLOR", "GETALPHA":
			return "ENTITY." + mn, true, true
		case "CREATEVEHICLE":
			return "BODY3D.CREATEVEHICLE", false, false
		case "SETTHROTTLE":
			return "VEHICLE.SETTHROTTLE", false, false
		case "SETSTEER":
			return "VEHICLE.SETSTEER", false, false
		}
	case heap.TagSprite:
		switch mn {
		case "SETPOS", "POS":
			return "SPRITE.GETPOS", true, true
		case "SETROT", "ROT":
			return "SPRITE.GETROT", true, true
		case "SETSCALE", "SCALE":
			return "SPRITE.GETSCALE", true, true
		case "SETCOLOR", "COLOR", "COL":
			return "SPRITE.GETCOLOR", true, true
		case "SETALPHA", "ALPHA", "A":
			return "SPRITE.GETALPHA", true, true
		}
	case heap.TagNavAgent:
		switch mn {
		case "SETPOS", "SETPOSITION", "POS", "POSITION":
			return "NAVAGENT.GETPOS", true, true
		case "SETROT", "SETROTATION", "ROTATION", "ROTATE", "ROT":
			return "NAVAGENT.GETROT", true, true
		case "SETSPEED", "SPEED":
			return "NAVAGENT.GETSPEED", true, true
		case "SETMAXFORCE", "MAXFORCE":
			return "NAVAGENT.GETMAXFORCE", true, true
		case "GETPOS", "GETROT", "GETSPEED", "GETMAXFORCE":
			return "NAVAGENT." + mn, true, true
		case "ISATDESTINATION":
			return "NAVAGENT.ISATDESTINATION", true, true
		}
	case heap.TagLight2D:
		switch mn {
		case "SETPOS", "SETPOSITION", "POS", "POSITION":
			return "LIGHT2D.GETPOS", true, true
		case "SETCOLOR", "COLOR", "COL":
			return "LIGHT2D.GETCOLOR", true, true
		case "SETRADIUS", "RADIUS":
			return "LIGHT2D.GETRADIUS", true, true
		case "SETINTENSITY", "INTENSITY":
			return "LIGHT2D.GETINTENSITY", true, true
		case "GETPOS", "GETCOLOR", "GETRADIUS", "GETINTENSITY":
			return "LIGHT2D." + mn, true, true
		}
	case heap.TagCamera2D:
		switch mn {
		case "SETPOS", "SETPOSITION", "POS", "POSITION":
			return "CAMERA2D.GETPOS", true, true
		case "SETTARGET", "TARGET":
			return "CAMERA2D.GETPOS", true, true
		case "SETROT", "SETROTATION", "ROTATION", "ROTATE", "ROT":
			return "CAMERA2D.GETROTATION", true, true
		case "GETPOS", "GETROTATION", "GETMATRIX", "GETZOOM", "GETOFFSET":
			return "CAMERA2D." + mn, true, true
		case "GETROT":
			return "CAMERA2D.GETROTATION", true, true
		case "SETZOOM", "ZOOM":
			return "CAMERA2D.GETZOOM", true, true
		case "SETOFFSET", "OFFSET":
			return "CAMERA2D.GETOFFSET", true, true
		}
	case heap.TagKinematicBody, heap.TagStaticBody, heap.TagTriggerBody:
		switch mn {
		case "SETPOS", "SETPOSITION", "POS", "POSITION":
			return "BODYREF.GETPOSITION", true, true
		case "SETROT", "SETROTATION", "ROTATION", "ROTATE", "ROT":
			return "BODYREF.GETROTATION", true, true
		case "GETPOSITION", "GETPOS":
			return "BODYREF.GETPOSITION", true, true
		case "GETROTATION", "GETROT":
			return "BODYREF.GETROTATION", true, true
		case "SETVELOCITY", "SETVEL", "GETVELOCITY", "GETVEL", "VELOCITY", "VEL":
			if tag == heap.TagKinematicBody {
				return "KINEMATICREF.GETVELOCITY", true, true
			}
		}
	case heap.TagTerrain:
		switch mn {
		case "SETPOS", "POS":
			return "TERRAIN.GETPOS", true, true
		case "SETSCALE", "SCALE":
			return "TERRAIN.GETSCALE", true, true
		case "SETDETAIL":
			return "TERRAIN.GETDETAIL", true, true
		}
	case heap.TagWater:
		switch mn {
		case "SETPOS", "POS":
			return "WATER.GETPOS", true, true
		case "SETCOLOR", "COLOR", "COL":
			return "WATER.GETCOLOR", true, true
		case "SETWAVEHEIGHT", "WAVEHEIGHT":
			return "WATER.GETWAVEHEIGHT", true, true
		// Builtin maps bare "WAVE" → WATER.SETWAVE; 0-arg should read speed (same as SETWAVE name).
		case "SETWAVE", "WAVE", "WAVESPEED":
			return "WATER.GETWAVESPEED", true, true
		case "SETSHALLOWCOLOR", "SHALLOWCOLOR":
			return "WATER.GETSHALLOWCOLOR", true, true
		case "SETDEEPCOLOR", "DEEPCOLOR":
			return "WATER.GETDEEPCOLOR", true, true
		}
	case heap.TagSky:
		switch mn {
		// Builtin maps setTime → SKY.SETTIME; 0-arg reads clock hours.
		case "SETTIME", "TIME", "GETTIMEHOURS":
			return "SKY.GETTIMEHOURS", true, true
		case "ISNIGHT", "NIGHT":
			return "SKY.ISNIGHT", true, true
		}
	case heap.TagCloud:
		switch mn {
		case "SETCOVERAGE", "COVERAGE", "GETCOVERAGE":
			return "CLOUD.GETCOVERAGE", true, true
		}
	case heap.TagWeather:
		switch mn {
		case "SETTYPE", "TYPE", "GETTYPE":
			return "WEATHER.GETTYPE", true, true
		case "SETCOVERAGE", "COVERAGE", "GETCOVERAGE":
			return "WEATHER.GETCOVERAGE", true, true
		}
	case heap.TagBiome:
		switch mn {
		case "SETTEMP", "TEMP", "GETTEMP":
			return "BIOME.GETTEMP", true, true
		case "SETHUMIDITY", "HUMIDITY", "GETHUMIDITY":
			return "BIOME.GETHUMIDITY", true, true
		}
	case heap.TagTable:
		switch mn {
		case "ROWCOUNT", "ROWS", "SETSIZE":
			return "TABLE.ROWCOUNT", true, true
		case "COLCOUNT", "COLS", "COLUMNS":
			return "TABLE.COLCOUNT", true, true
		}
	case heap.TagJSON:
		switch mn {
		case "LEN", "LENGTH", "COUNT", "SETSIZE":
			return "JSON.LEN", true, true
		case "TYPE":
			return "JSON.TYPE", true, true
		case "KEYS":
			return "JSON.KEYS", true, true
		}
	case heap.TagCSV:
		switch mn {
		case "ROWCOUNT", "ROWS", "SETSIZE":
			return "CSV.ROWCOUNT", true, true
		case "COLCOUNT", "COLS", "COLUMNS":
			return "CSV.COLCOUNT", true, true
		}
	case heap.TagMem:
		switch mn {
		// "size" normalizes to SETSIZE; byte-length getter uses MEM.SIZE.
		case "SIZE", "SETSIZE", "LENGTH", "LEN":
			return "MEM.SIZE", true, true
		}
	case heap.TagLobby:
		switch mn {
		case "GETNAME", "NAME":
			return "LOBBY.GETNAME", true, true
		}
	case heap.TagNetPacket:
		switch mn {
		case "DATA":
			return "PACKET.DATA", true, true
		}
	case heap.TagDB:
		switch mn {
		case "ISOPEN":
			return "DB.ISOPEN", true, true
		}
	case heap.TagDecal:
		switch mn {
		case "SETPOS", "POS":
			return "DECAL.GETPOS", true, true
		case "SETSIZE", "SIZE", "GETSIZE":
			return "DECAL.GETSIZE", true, true
		case "SETLIFETIME", "LIFETIME", "GETLIFETIME":
			return "DECAL.GETLIFETIME", true, true
		case "SETROT", "ROT", "GETROT":
			return "DECAL.GETROT", true, true
		}
	case heap.TagTexture, heap.TagImage, heap.TagImageSequence:
		pre := "TEXTURE"
		if tag == heap.TagImage || tag == heap.TagImageSequence {
			pre = "IMAGE"
		}
		switch mn {
		// "size" normalizes to SETSIZE; route 0-arg read to GETSIZE (same as WIDTH/HEIGHT).
		case "SETSIZE", "SIZE", "GETSIZE":
			return pre + ".GETSIZE", true, true
		case "GETWIDTH", "WIDTH":
			return pre + ".GETWIDTH", true, true
		case "GETHEIGHT", "HEIGHT":
			return pre + ".GETHEIGHT", true, true
		}
	case heap.TagTilemap:
		switch mn {
		case "WIDTH", "GETWIDTH":
			return "TILEMAP.WIDTH", true, true
		case "HEIGHT", "GETHEIGHT":
			return "TILEMAP.HEIGHT", true, true
		// SETSIZE from ".size()": no single "area" builtin; width in tiles is the usual scalar.
		case "SETSIZE":
			return "TILEMAP.WIDTH", true, true
		}
	case heap.TagPool:
		// POOL.GET / POOL.FREE take the pool handle as the only argument.
		switch mn {
		case "GET":
			return "POOL.GET", true, true
		case "FREE":
			return "POOL.FREE", true, true
		}
	case heap.TagAtlas:
		// ATLAS.FREE is (atlas). GETSPRITE prepends atlas; name args follow when present.
		switch mn {
		case "GETSPRITE":
			return "ATLAS.GETSPRITE", true, true
		case "FREE":
			return "ATLAS.FREE", true, true
		}
	case heap.TagNav:
		switch mn {
		case "FREE":
			return "NAV.FREE", true, true
		}
	case heap.TagPath:
		switch mn {
		case "ISVALID", "VALID":
			return "PATH.ISVALID", true, true
		case "NODECOUNT", "COUNT", "LEN", "SETSIZE":
			return "PATH.NODECOUNT", true, true
		case "FREE":
			return "PATH.FREE", true, true
		}
	case heap.TagFile:
		switch mn {
		// normalize maps "pos"/"position" → SETPOS; 0-arg must read (FILE.SEEK is for seek/set).
		case "GETPOS", "SETPOS", "POS", "TELL":
			return "FILE.GETPOS", true, true
		case "GETSIZE", "SIZE", "SETSIZE":
			return "FILE.GETSIZE", true, true
		case "GETEOF", "EOF":
			return "FILE.GETEOF", true, true
		}
	case heap.TagArray:
		switch mn {
		case "GETSIZE", "GETLEN", "LEN", "COUNT", "SIZE", "SETSIZE":
			return "ARRAY.LEN", true, true
		}
	case heap.TagSound, heap.TagMusic:
		switch mn {
		// normalize maps setVolume → SETVOLUME etc.; 0-arg must read, not AUDIO.SET*VOLUME.
		case "GETVOLUME", "VOLUME", "SETVOLUME":
			if tag == heap.TagSound {
				return "AUDIO.GETSOUNDVOLUME", true, true
			}
			return "AUDIO.GETMUSICVOLUME", true, true
		case "GETPITCH", "PITCH", "SETPITCH":
			if tag == heap.TagSound {
				return "AUDIO.GETSOUNDPITCH", true, true
			}
			return "AUDIO.GETMUSICPITCH", true, true
		case "GETPAN", "PAN", "SETPAN":
			if tag != heap.TagSound {
				break
			}
			return "AUDIO.GETSOUNDPAN", true, true
		case "GETLENGTH", "LENGTH":
			if tag != heap.TagMusic {
				break
			}
			return "AUDIO.GETMUSICLENGTH", true, true
		case "GETTIME", "TIME":
			if tag != heap.TagMusic {
				break
			}
			return "AUDIO.GETMUSICTIME", true, true
		}
	case heap.TagHost:
		switch mn {
		case "UPDATE", "PUMP":
			return "NET.UPDATE", true, true
		case "RECEIVE", "POP":
			return "NET.RECEIVE", true, true
		case "PEERCOUNT", "PEERS":
			return "NET.PEERCOUNT", true, true
		case "CLOSE":
			return "NET.CLOSE", true, true
		}
	case heap.TagEvent:
		switch mn {
		case "TYPE":
			return "EVENT.TYPE", true, true
		case "PEER":
			return "EVENT.PEER", true, true
		case "DATA":
			return "EVENT.DATA", true, true
		case "CHANNEL":
			return "EVENT.CHANNEL", true, true
		case "FREE":
			return "EVENT.FREE", true, true
		}
	case heap.TagAutomationList:
		switch mn {
		case "LISTCOUNT", "SETSIZE":
			return "EVENT.LISTCOUNT", true, true
		}
	case heap.TagPlayer2D:
		switch mn {
		// "pos"/"position" normalize to SETPOS; 0-arg must read XZ, not PLAYER2D.SETPOS.
		case "SETPOS", "POS", "POSITION":
			return "PLAYER2D.GETPOS", true, true
		case "GETPOS":
			return "PLAYER2D.GETPOS", true, true
		case "GETX", "X":
			return "PLAYER2D.GETX", true, true
		case "GETZ", "Z":
			return "PLAYER2D.GETZ", true, true
		}
	case heap.TagGameTimer:
		switch mn {
		case "FINISHED", "DONE":
			return "TIMER.FINISHED", true, true
		case "REMAINING", "LEFT":
			return "TIMER.REMAINING", true, true
		case "FREE":
			return "TIMER.FREE", true, true
		}
	case heap.TagGameTimerSim:
		switch mn {
		case "DONE":
			return "TIMER.DONE", true, true
		case "FRACTION", "FRAC":
			return "TIMER.FRACTION", true, true
		case "REMAINING", "LEFT":
			return "TIMER.REMAINING", true, true
		case "SETLOOP", "LOOP", "GETLOOP":
			return "TIMER.GETLOOP", true, true
		case "FREE":
			return "TIMER.FREE", true, true
		}
	case heap.TagGameStopwatch:
		switch mn {
		case "RESET":
			return "STOPWATCH.RESET", true, true
		case "ELAPSED", "TIME":
			return "STOPWATCH.ELAPSED", true, true
		case "FREE":
			return "STOPWATCH.FREE", true, true
		}
	case heap.TagVec2:
		switch mn {
		case "X":
			return "VEC2.X", true, true
		case "Y":
			return "VEC2.Y", true, true
		case "LEN", "LENGTH":
			return "VEC2.LENGTH", true, true
		}
	case heap.TagVec3:
		switch mn {
		case "X":
			return "VEC3.X", true, true
		case "Y":
			return "VEC3.Y", true, true
		case "Z":
			return "VEC3.Z", true, true
		case "LEN", "LENGTH":
			return "VEC3.LENGTH", true, true
		}
	case heap.TagColor:
		switch mn {
		case "R":
			return "COLOR.R", true, true
		case "G":
			return "COLOR.G", true, true
		case "B":
			return "COLOR.B", true, true
		case "SETALPHA", "ALPHA":
			return "COLOR.A", true, true
		}
	case heap.TagAudioStream:
		switch mn {
		case "ISPLAYING", "PLAYING":
			return "AUDIOSTREAM.ISPLAYING", true, true
		case "READY", "ISREADY":
			return "AUDIOSTREAM.ISREADY", true, true
		case "GETVOLUME", "SETVOLUME", "VOLUME":
			return "AUDIOSTREAM.GETVOLUME", true, true
		case "GETPITCH", "SETPITCH", "PITCH":
			return "AUDIOSTREAM.GETPITCH", true, true
		case "GETPAN", "SETPAN", "PAN":
			return "AUDIOSTREAM.GETPAN", true, true
		}
	}
	return handleCallBuiltin(tag, method)
}

// HandleCallSuggestions lists common script-side method names for a handle type (error hints).
func HandleCallSuggestions(tag uint16) []string {
	var out []string
	switch tag {
	case heap.TagCamera:
		out = []string{"Begin", "End", "FOV", "Free", "GetMatrix", "GetPos", "GetRay", "GetTarget", "GetViewRay", "GetYaw", "IsOnScreen",
			"Look", "LookAt", "MouseRay", "Move", "Orbit", "Pos", "SetFOV", "SetOrbit", "SetPos", "SetPosition", "SetProjection", "SetTarget", "SetUp", "WorldToScreen", "Yaw", "Zoom"}
	case heap.TagEntityRef:
		out = []string{"A", "Col", "Color", "Free", "Hide", "LockAxis", "Move", "MoveWithCamera", "Pos", "Rot", "Scale", "SetBounciness", "SetCCD", "SetCollisionMesh", "SetDamping", "SetGravityFactor", "SetStatic", "Show", "Turn"}
	case heap.TagCamera2D:
		out = []string{"Begin", "End", "Free", "GetMatrix", "GetOffset", "GetPos", "GetRotation", "GetZoom", "ScreenToWorld", "SetOffset", "SetRotation", "SetTarget", "SetZoom", "WorldToScreen"}
	case heap.TagRenderTexture:
		out = []string{"Begin", "End", "Free", "Texture"}
	case heap.TagSky:
		out = []string{"Draw", "Free", "GetTimeHours", "IsNight", "SetDayLength", "SetTime", "Time", "Update"}
	case heap.TagCloud:
		out = []string{"Coverage", "Draw", "Free", "GetCoverage", "SetCoverage", "Update"}
	case heap.TagWeather:
		out = []string{"Coverage", "Draw", "Free", "GetCoverage", "GetType", "SetCoverage", "SetType", "Type", "Update"}
	case heap.TagScatterSet:
		out = []string{"Apply", "DrawAll", "Free"}
	case heap.TagProp:
		out = []string{"Free", "Place"}
	case heap.TagBiome:
		out = []string{"Free", "GetHumidity", "GetTemp", "Humidity", "SetHumidity", "SetTemp", "Temp"}
	case heap.TagNoise:
		out = []string{"FillArray", "FillArrayNorm", "FillImage", "Free", "Get", "Get3D", "GetDomainWarped", "GetNorm", "GetTileable",
			"SetCellularDistance", "SetCellularJitter", "SetCellularType", "SetDomainWarpAmplitude", "SetDomainWarpType", "SetFrequency", "SetGain", "SetLacunarity", "SetOctaves", "SetPingPongStrength", "SetSeed", "SetType", "SetWeightedStrength"}
	case heap.TagTable:
		out = []string{"AddRow", "ColCount", "Free", "Get", "RowCount", "Set", "ToCSV", "ToJSON"}
	case heap.TagPool:
		out = []string{"Free", "Get", "Prewarm", "Return", "SetFactory", "SetReset"}
	case heap.TagJSON:
		out = []string{"Append", "Clear", "Delete", "Free", "GetArray", "GetBool", "GetFloat", "GetInt", "GetObject", "GetString", "Has", "Keys", "Len", "Minify", "Pretty", "Query", "SaveFile", "SetBool", "SetFloat", "SetInt", "SetNull", "SetString", "ToCSV", "ToFile", "ToFilePretty", "ToString", "Type"}
	case heap.TagCSV:
		out = []string{"ColCount", "Free", "Get", "RowCount", "Save", "Set", "ToJSON", "ToString"}
	case heap.TagDB:
		out = []string{"Begin", "Changes", "Close", "Exec", "IsOpen", "LastInsertID", "Prepare", "Query", "QueryJSON"}
	case heap.TagDBRows:
		out = []string{"Close", "GetFloat", "GetInt", "GetString", "Next"}
	case heap.TagDBStmt:
		out = []string{"Close", "Exec", "StmtClose", "StmtExec"}
	case heap.TagDBTx:
		out = []string{"Commit", "Rollback"}
	case heap.TagRng:
		out = []string{"Free", "Next", "NextF"}
	case heap.TagMem:
		out = []string{"Clear", "Free", "GetByte", "GetDouble", "GetDword", "GetFloat", "GetString", "GetWord", "Resize", "SetByte", "SetDouble", "SetDword", "SetFloat", "SetString", "SetWord", "Size"}
	case heap.TagLobby:
		out = []string{"Free", "GetName", "Join", "SetHost", "SetProperty", "Start"}
	case heap.TagNetPacket:
		out = []string{"Data", "Free"}
	case heap.TagHost:
		out = []string{"Broadcast", "Close", "Flush", "PeerCount", "Receive", "Service", "SetBandwidth", "Update"}
	case heap.TagEvent:
		out = []string{"Channel", "Data", "Free", "Peer", "Type"}
	case heap.TagPlayer2D:
		out = []string{"Clamp", "Free", "GetPos", "GetX", "GetZ", "KeepInBounds", "Move", "SetPos"}
	case heap.TagGameTimer:
		out = []string{"Finished", "Free", "Remaining", "Reset"}
	case heap.TagGameTimerSim:
		out = []string{"Done", "Fraction", "Free", "GetLoop", "Loop", "Remaining", "Rewind", "SetLoop", "Start", "Stop", "Update"}
	case heap.TagGameStopwatch:
		out = []string{"Elapsed", "Free", "Reset"}
	case heap.TagDecal:
		out = []string{"Draw", "Free", "GetLifetime", "GetPos", "GetRot", "GetSize", "Lifetime", "Pos", "Rot", "SetLifetime", "SetPos", "SetSize", "Size"}
	case heap.TagTilemap:
		out = []string{"CollisionAt", "Draw", "DrawLayer", "Free", "GetTile", "Height", "IsSolid", "IsSolidCategory",
			"LayerCount", "LayerName", "MergeCollisionLayer", "SetCollision", "SetTile", "SetTileSize", "Width"}
	case heap.TagAtlas:
		out = []string{"Free", "GetSprite"}
	case heap.TagLight2D:
		out = []string{"Free", "GetColor", "GetIntensity", "GetPos", "GetRadius", "SetColor", "SetIntensity", "SetPos", "SetRadius"}
	case heap.TagPhysicsBody:
		out = []string{"Activate", "AddForce", "AddImpulse", "ApplyForce", "ApplyImpulse", "BufferIndex", "Collided", "CollisionNormal", "CollisionOther", "CollisionPoint", "Deactivate", "Force", "Free", "GetLinearVel", "Impulse", "Pos", "Rot", "Scale", "SetLinearVel", "SetPos", "SetPosition", "SetRot", "SetVelocity", "Vel", "Velocity"}
	case heap.TagPhysicsBuilder:
		out = []string{"AddBox", "AddCapsule", "AddMesh", "AddSphere", "Commit", "Free"}
	case heap.TagBody2D:
		out = []string{"AddCircle", "AddPolygon", "AddRect", "ApplyForce", "ApplyImpulse", "Collided", "CollisionNormal", "CollisionOther", "CollisionPoint", "Commit", "Free", "GetMass", "GetRestitution", "GetVelocity", "Pos", "SetFriction", "SetMass", "SetPos", "SetRestitution", "SetVel", "SetVelocity", "Vel", "Velocity"}
	case heap.TagPeer:
		out = []string{"Disconnect", "IP", "Ping", "Send", "SendPacket"}
	case heap.TagFile:
		out = []string{"Close", "EOF", "Free", "GetEOF", "GetPos", "GetSize", "Pos", "Read", "ReadBuf", "ReadLine", "ReadString", "Seek", "Size", "Tell", "Write", "WriteLn"}
	case heap.TagArray:
		out = []string{"Contains", "Copy", "Count", "Fill", "Find", "Free", "Join", "Len", "Pop", "Push", "Reverse", "Shift", "Size", "Slice", "Sort", "Splice", "Unshift"}
	case heap.TagSound:
		out = []string{"Free", "Pause", "Play", "Resume", "SetVolume", "Stop", "Volume"}
	case heap.TagMusic:
		out = []string{"Free", "Length", "Pause", "Play", "Resume", "Seek", "SetPitch", "SetVolume", "Stop", "Time", "Update", "Volume"}
	case heap.TagAudioStream:
		out = []string{"Free", "GetPan", "GetPitch", "GetVolume", "IsPlaying", "IsReady", "Pan", "Pause", "Pitch", "Play", "Resume", "SetPan", "SetPitch", "SetVolume", "Stop", "Update", "Volume"}
	case heap.TagWave:
		out = []string{"Copy", "Crop", "Export", "Format", "Free"}
	case heap.TagCharController:
		out = []string{"AddVel", "AddVelocity", "AirControl", "DrainContacts", "Free", "Friction", "GetAirControl", "GetBounce", "GetCeiling", "GetFriction", "GetGravity", "GetGroundControl", "GetGroundVelocity", "GetIsSliding", "GetJumpBuffer", "GetMaxSlope", "GetPadding", "GetPos", "GetPosition", "GetSlopeAngle", "GetSnapDistance", "GetSpeed", "GetStepHeight", "Gravity", "GroundControl", "Grounded", "IsGrounded", "Jump", "JumpBuffer", "MaxSlope", "Move", "MoveWithCamera", "OnSlope", "OnWall", "Padding", "Pos", "SetAirControl", "SetBounce", "SetContactListener", "SetFriction", "SetGravity", "SetGravityScale", "SetGroundControl", "SetJumpBuffer", "SetLinearVel", "SetLinearVelocity", "SetMaxSlope", "SetPadding", "SetPos", "SetPosition", "SetSetting", "SetSnapDistance", "SetStepHeight", "SetStickDown", "SetVelocity", "SlopeAngle", "Snap", "SnapDistance", "Speed", "StepHeight", "Update", "UpdateMove"}
	case heap.TagTween:
		out = []string{"Free", "GetLoop", "GetYoyo", "IsFinished", "IsPlaying", "Loop", "OnComplete", "Progress", "Start", "Stop", "Then", "To", "Update", "Yoyo"}
	case heap.TagRay:
		out = []string{"Free", "GetDir", "GetPos", "SetDir", "SetPos"}
	case heap.TagBBox:
		out = []string{"Check", "CheckSphere", "Free", "GetMax", "GetMin", "SetMax", "SetMin"}
	case heap.TagBSphere:
		out = []string{"Check", "CheckBox", "Free", "GetPos", "GetRadius", "SetPos", "SetRadius"}
	case heap.TagShape:
		out = []string{"Free"}
	case heap.TagKinematicBody:
		out = []string{"EnableCollision", "Free", "Pos", "Rot", "SetLayer", "SetPos", "SetPosition", "SetRot", "SetRotation", "SetVel", "SetVelocity", "Update"}
	case heap.TagStaticBody, heap.TagTriggerBody:
		out = []string{"EnableCollision", "Free", "Pos", "Rot", "SetLayer", "SetPos", "SetPosition", "SetRot", "SetRotation"}
	case heap.TagSprite:
		out = []string{"Alpha", "Col", "Color", "Draw", "Free", "Pos", "Rot", "Scale", "SetAlpha", "SetColor", "SetOrigin", "SetPos", "SetPosition", "SetRot", "SetScale", "SetFrame", "DefAnim", "PlayAnim", "UpdateAnim", "Play", "Hit", "Collide", "PointHit"}
	case heap.TagModel, heap.TagLODModel:
		out = []string{"AttachTo", "Clone", "Cull", "Depth", "Detach", "Draw", "Diffuse", "Emissive", "Exists", "Fog", "Free", "GetPos", "GetRot", "GetScale",
			"Lighting", "Metal", "Move", "Pos", "Rotate", "Rot", "Rough", "Scale", "ScrollTexture", "SetAlpha", "SetAmbientColor", "SetBlend", "SetColor", "SetCull", "SetDepth", "SetFog", "SetGPUSkinning", "SetLighting", "SetLoDDistances", "SetMaterial", "SetMetal", "SetPos", "SetRough", "SetRot", "SetScale", "SetSpecularPow", "SetStageBlend", "SetStageRotate", "SetStageScale", "SetStageScroll", "SetTextureStage", "SetWireframe", "Specular", "Wireframe"}
	case heap.TagMaterial:
		out = []string{"Col", "Color", "Free", "SetColor", "SetEffect", "SetEffectParam", "SetFloat", "SetShader", "SetTexture"}
	case heap.TagShader:
		out = []string{"Free", "GetLoc", "SetFloat", "SetInt", "SetTexture", "SetVec2", "SetVec3", "SetVec4"}
	case heap.TagParticle:
		out = []string{"Count", "Draw", "Free", "GetAlpha", "GetColor", "GetPos", "GetSize", "GetVelocity", "IsAlive", "Play", "Scale", "SetBillboard", "SetBurst", "SetColor", "SetColorEnd",
			"SetDirection", "SetEmitRate", "SetEndColor", "SetEndSize", "SetGravity", "SetLifetime", "SetPos", "SetPosition",
			"SetRate", "SetSize", "SetSpeed", "SetSpread", "SetStartColor", "SetStartSize", "SetTexture", "SetVelocity", "Stop", "Update", "Vel"}
	case heap.TagTerrain:
		out = []string{"ApplyMap", "ApplyTiles", "Detail", "Draw", "FillFlat", "FillPerlin", "Free", "GetDetail", "GetHeight", "GetNormal", "GetPos", "GetScale", "GetSlope", "GetSplat", "Load", "Lower", "Place", "Pos", "Raise", "Raycast", "Scale", "SetAsyncMeshBuild", "SetChunkSize", "SetDetail", "SetMeshBuildBudget", "SetPos", "SetScale", "SnapY"}
	case heap.TagWater:
		out = []string{"AutoPhysics", "Draw", "Free", "GetColor", "GetDeepColor", "GetDepth", "GetPos", "GetShallowColor", "GetWaveHeight", "GetWaveSpeed", "GetWaveY", "IsUnder", "Pos", "SetColor", "SetDeepColor", "SetPos", "SetShallowColor", "SetWave", "SetWaveHeight", "Update"}
	case heap.TagInstancedModel:
		out = []string{"Count", "Draw", "DrawLOD", "Free", "GetAlpha", "GetColor", "GetPos", "GetRot", "GetScale", "SetColor", "SetCullDistance", "SetInstancePos", "SetInstanceScale",
			"SetMatrix", "SetPos", "SetRot", "SetScale", "UpdateBuffer", "UpdateInstances"}
	case heap.TagLight:
		out = []string{"Enable", "Free", "IsEnabled", "SetColor", "SetDir", "SetInnerCone", "SetIntensity",
			"SetOuterCone", "SetPos", "SetPosition", "SetRange", "SetShadow", "SetShadowBias", "SetTarget"}
	case heap.TagMatrix:
		out = []string{"ApplyX", "ApplyY", "ApplyZ", "Free", "FromMat4", "GetElement", "Inverse", "Multiply", "SetRotation", "ToQuat", "Transpose"}
	case heap.TagMesh:
		out = []string{"Draw", "DrawRotated", "Free", "TriangleCount", "VertexCount"}
	case heap.TagFont:
		out = []string{"DrawDefault", "Free", "SetDefault"}
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
	case heap.TagNav:
		out = []string{"AddObstacle", "AddTerrain", "Build", "DebugDraw", "FindPath", "Free", "SetGrid"}
	case heap.TagPath:
		out = []string{"Free", "IsValid", "NodeCount", "NodeX", "NodeY", "NodeZ"}
	case heap.TagNavAgent:
		out = []string{"ApplyForce", "Free", "GetMaxForce", "GetPos", "GetRot", "GetSpeed", "IsAtDestination", "MaxForce", "MoveTo", "Pos", "Rot", "SetMaxForce", "SetPos", "SetPosition", "SetSpeed", "Stop", "Update", "X", "Y", "Z"}
	case heap.TagSteerGroup:
		out = []string{"Add", "Clear"}
	case heap.TagBTree:
		out = []string{"AddAction", "AddCondition", "Free", "Run", "Sequence"}
	case heap.TagComputeShader:
		out = []string{"Dispatch", "Free", "SetBuffer", "SetFloat", "SetInt"}
	case heap.TagShaderBuffer:
		out = []string{"Free"}
	case heap.TagJoint2D:
		out = []string{"Free"}
	case heap.TagBrush:
		out = []string{"Alpha", "Blend", "Color", "Free", "FX", "Shininess", "Texture"}
	case heap.TagMeshBuilder:
		out = []string{"AddTriangle", "AddVertex", "VertexX", "VertexY", "VertexZ", "X", "Y", "Z"}
	case heap.TagImageSequence:
		out = []string{"Free", "Height", "Size", "Width"}
	case heap.TagSpriteGroup:
		out = []string{"Add", "Clear", "Draw", "Free", "Remove"}
	case heap.TagSpriteLayer:
		out = []string{"Add", "Clear", "Draw", "Free", "SetZ"}
	case heap.TagSpriteBatch:
		out = []string{"Add", "Clear", "Draw", "Free"}
	case heap.TagSpriteUI:
		out = []string{"Draw", "Free"}
	case heap.TagParticle2D:
		out = []string{"Draw", "Emit", "Free", "Update"}
	case heap.TagQuaternion:
		out = []string{"Free", "Invert", "Multiply", "Normalize", "Slerp", "ToEuler", "ToMat4", "Transform"}
	case heap.TagColor:
		out = []string{"Brightness", "Contrast", "Fade", "Free", "Invert", "Lerp", "ToHex", "ToHSV", "ToHSVX", "ToHSVY", "ToHSVZ"}
	case heap.TagVec2:
		out = []string{"Add", "Angle", "Dist", "Distance", "Free", "Len", "Length", "Lerp", "Mul", "Normalize", "Rotate", "Set", "Sub", "TransformMat4", "X", "Y"}
	case heap.TagVec3:
		out = []string{"Add", "Angle", "Cross", "Dist", "Distance", "Dot", "Equals", "Free", "Len", "Length", "Lerp", "Mul", "Negate", "Normalize", "OrthoNormalize", "Project", "Reflect", "RotateByQuat", "Set", "Sub", "TransformMat4", "X", "Y", "Z"}
	case heap.TagAutomationList:
		out = []string{"Export", "ListClear", "ListCount", "ListFree", "Replay", "SetActiveList"}
	case heap.TagTacticalGrid:
		out = []string{"Draw", "FollowTerrain", "Free", "GetCell", "GetNeighbors", "GetPath", "PlaceEntity", "Raycast", "SetCell", "Snap", "WorldToCell"}
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
	return msg + "\n  Hint: See docs/API_consistency.md for handle methods vs NS.COMMAND calls."
}
