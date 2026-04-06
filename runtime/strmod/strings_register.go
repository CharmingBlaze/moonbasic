package strmod

import "moonbasic/runtime"

// registerStringBuiltins registers STR$/INT, slice/search/format helpers, SPLIT$/JOIN$, etc.
func registerStringBuiltins(r runtime.Registrar) {
	registerStringsConv(r)
	registerStringsSlice(r)
	registerStringsSearch(r)
	registerStringsFormat(r)
	registerStringsCheck(r)
	registerStringsSplitJoin(r)
}
