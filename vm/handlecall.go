package vm

import (
	"sort"
	"strings"

	"moonbasic/runtime"
)

// normalizeHandleMethod maps legacy or verbose names to the canonical verb used in the registry.
func normalizeHandleMethod(mn string) string {
	switch mn {
	case "SETPOSITION":
		return "SETPOS"
	default:
		return mn
	}
}

// handleCallRegistryPrefix returns the dotted namespace prefix for handle-typed registry keys (e.g. CAMERA for Camera3D).
func handleCallRegistryPrefix(typeName string) string {
	switch strings.ToUpper(strings.TrimSpace(typeName)) {
	case "CAMERA3D":
		return "CAMERA."
	case "CAMERA2D":
		return "CAMERA2D."
	case "MATRIX4":
		return "TRANSFORM."
	default:
		return strings.ToUpper(strings.TrimSpace(typeName)) + "."
	}
}

// handleCallBuiltin maps heap TypeName + script method to a registry command key and whether
// the receiver handle is passed as the first argument to that builtin.
func handleCallBuiltin(typeName, method string) (registryKey string, prependReceiver bool, ok bool) {
	tn := strings.ToUpper(strings.TrimSpace(typeName))
	mn := normalizeHandleMethod(strings.ToUpper(strings.TrimSpace(method)))
	switch tn {
	case "CAMERA3D":
		switch mn {
		case "END":
			return "CAMERA.END", false, true
		case "BEGIN", "SETPOS", "SETTARGET", "SETFOV", "MOVE", "GETRAY", "GETVIEWRAY", "GETMATRIX":
			return "CAMERA." + mn, true, true
		}
	case "CAMERA2D":
		switch mn {
		case "END":
			return "CAMERA2D.END", false, true
		case "BEGIN", "SETTARGET", "SETOFFSET", "SETZOOM", "SETROTATION":
			return "CAMERA2D." + mn, true, true
		}
	case "TILEMAP":
		switch mn {
		case "FREE", "SETTILESIZE", "DRAW", "DRAWLAYER", "GETTILE", "SETTILE", "ISSOLID", "ISSOLIDCATEGORY",
			"WIDTH", "HEIGHT", "LAYERCOUNT", "LAYERNAME", "COLLISIONAT", "SETCOLLISION", "MERGECOLLISIONLAYER":
			return "TILEMAP." + mn, true, true
		}
	case "ATLAS":
		switch mn {
		case "FREE", "GETSPRITE":
			return "ATLAS." + mn, true, true
		}
	case "LIGHT2D":
		switch mn {
		case "FREE", "SETPOS", "SETCOLOR", "SETRADIUS", "SETINTENSITY":
			return "LIGHT2D." + mn, true, true
		}
	case "BODY3D":
		if mn == "SETPOS" {
			return "BODY3D.SETPOS", true, true
		}
	case "CHARCONTROLLER":
		if mn == "SETPOS" {
			return "CHARCONTROLLER.SETPOS", true, true
		}
	case "SPRITE":
		switch mn {
		case "SETPOS", "DRAW", "DEFANIM", "PLAYANIM", "UPDATEANIM", "HIT":
			return "SPRITE." + mn, true, true
		}
	case "LIGHT":
		switch mn {
		case "SETDIR", "SETSHADOW":
			return "LIGHT." + mn, true, true
		}
	case "MODEL":
		switch mn {
		case "SETPOS":
			return "MODEL.SETPOS", true, true
		case "DRAW":
			return "MODEL.DRAW", true, true
		}
	case "LODMODEL":
		switch mn {
		case "SETPOS", "SETPOSITION":
			return "MODEL.SETPOS", true, true
		case "DRAW":
			return "MODEL.DRAW", true, true
		}
	case "INSTANCEDMODEL":
		switch mn {
		case "DRAW":
			return "MODEL.DRAW", true, true
		case "SETINSTANCEPOS":
			return "MODEL.SETINSTANCEPOS", true, true
		case "SETINSTANCESCALE":
			return "MODEL.SETINSTANCESCALE", true, true
		case "UPDATEINSTANCES":
			return "MODEL.UPDATEINSTANCES", true, true
		}
	case "PARTICLE":
		switch mn {
		case "SETTEXTURE":
			return "PARTICLE.SETTEXTURE", true, true
		case "SETEMITRATE":
			return "PARTICLE.SETEMITRATE", true, true
		case "SETLIFETIME":
			return "PARTICLE.SETLIFETIME", true, true
		case "SETVELOCITY":
			return "PARTICLE.SETVELOCITY", true, true
		case "SETCOLOR":
			return "PARTICLE.SETCOLOR", true, true
		case "SETCOLOREND":
			return "PARTICLE.SETCOLOREND", true, true
		case "SETSIZE":
			return "PARTICLE.SETSIZE", true, true
		case "SETGRAVITY":
			return "PARTICLE.SETGRAVITY", true, true
		case "SETPOS":
			return "PARTICLE.SETPOS", true, true
		case "PLAY":
			return "PARTICLE.PLAY", true, true
		case "UPDATE":
			return "PARTICLE.UPDATE", true, true
		case "DRAW":
			return "PARTICLE.DRAW", true, true
		case "FREE":
			return "PARTICLE.FREE", true, true
		}
	case "MATRIX4":
		if mn == "SETROTATION" {
			return "TRANSFORM.SETROTATION", true, true
		}
	case "MESH":
		switch mn {
		case "DRAW", "DRAWROTATED":
			return "MESH." + mn, true, true
		}
	}
	return "", false, false
}

// HandleCallSuggestions lists common script-side method names for a handle TypeName (error hints).
func HandleCallSuggestions(typeName string) []string {
	tn := strings.ToUpper(strings.TrimSpace(typeName))
	var out []string
	switch tn {
	case "CAMERA3D":
		out = []string{"Begin", "End", "SetPos", "SetPosition", "SetTarget", "SetFOV", "Move", "GetRay", "GetViewRay", "GetMatrix"}
	case "CAMERA2D":
		out = []string{"Begin", "End", "SetOffset", "SetRotation", "SetTarget", "SetZoom"}
	case "TILEMAP":
		out = []string{"CollisionAt", "Draw", "DrawLayer", "Free", "GetTile", "Height", "IsSolid", "IsSolidCategory",
			"LayerCount", "LayerName", "MergeCollisionLayer", "SetCollision", "SetTile", "SetTileSize", "Width"}
	case "ATLAS":
		out = []string{"Free", "GetSprite"}
	case "LIGHT2D":
		out = []string{"Free", "SetColor", "SetIntensity", "SetPos", "SetRadius"}
	case "BODY3D":
		out = []string{"SetPos", "SetPosition"}
	case "CHARCONTROLLER":
		out = []string{"SetPos", "SetPosition"}
	case "SPRITE":
		out = []string{"Draw", "SetPos", "SetPosition", "DefAnim", "PlayAnim", "UpdateAnim", "Hit"}
	case "MODEL":
		out = []string{"Draw", "SetPos", "SetPosition"}
	case "LODMODEL":
		out = []string{"Draw", "SetPos", "SetPosition"}
	case "PARTICLE":
		out = []string{"Draw", "Free", "Play", "SetColor", "SetColorEnd", "SetEmitRate", "SetGravity", "SetLifetime", "SetPos", "SetPosition", "SetSize", "SetTexture", "SetVelocity", "Update"}
	case "INSTANCEDMODEL":
		out = []string{"Draw", "SetInstancePos", "SetInstanceScale", "UpdateInstances"}
	case "LIGHT":
		out = []string{"SetDir", "SetShadow"}
	case "MATRIX4":
		out = []string{"SetRotation"}
	case "MESH":
		out = []string{"Draw", "DrawRotated"}
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
func (v *VM) formatHandleCallError(typeName, methodName, callKey string, mapped bool, err error) string {
	msg := err.Error()
	if mapped {
		return msg
	}
	prefix := handleCallRegistryPrefix(typeName)
	keys := v.Registry.CommandKeys()
	prefixed := filterRegistryKeysByPrefix(keys, prefix)
	if alt, ok := runtime.BestSimilarCommand(callKey, prefixed, 3); ok {
		return msg + "\n  Did you mean " + alt + "?"
	}
	if sug := HandleCallSuggestions(typeName); len(sug) > 0 {
		return msg + "\n  Hint: For " + typeName + " handles use methods like " + strings.Join(sug, ", ") + "."
	}
	return msg + "\n  Hint: See docs/API_CONSISTENCY.md for handle methods vs NS.COMMAND calls."
}
