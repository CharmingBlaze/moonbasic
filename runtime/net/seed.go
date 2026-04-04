package mbnet

import "moonbasic/vm/value"

// SeedMultiplayerGlobals installs SYNC_TRANSFORM / SYNC_ANIMATION bit flags for SERVER.SYNCENTITY.
func SeedMultiplayerGlobals(globals map[string]value.Value) {
	globals["SYNC_TRANSFORM"] = value.FromFloat(1)
	globals["SYNC_ANIMATION"] = value.FromFloat(2)
}
