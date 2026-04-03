package vm

import "strings"

// handleCallBuiltin maps heap TypeName + script method to a registry command key and whether
// the receiver handle is passed as the first argument to that builtin.
func handleCallBuiltin(typeName, method string) (registryKey string, prependReceiver bool, ok bool) {
	tn := strings.ToUpper(strings.TrimSpace(typeName))
	mn := strings.ToUpper(strings.TrimSpace(method))
	switch tn {
	case "CAMERA3D":
		switch mn {
		case "END":
			return "CAMERA.END", false, true
		case "BEGIN", "SETPOS", "SETTARGET", "SETFOV", "MOVE", "GETRAY", "GETMATRIX":
			return "CAMERA." + mn, true, true
		}
	case "MATRIX4":
		switch mn {
		case "SETROTATION":
			return "MAT4.SETROTATION", true, true
		}
	case "MESH":
		switch mn {
		case "DRAW", "DRAWROTATED":
			return "MESH." + mn, true, true
		}
	}
	return "", false, false
}
