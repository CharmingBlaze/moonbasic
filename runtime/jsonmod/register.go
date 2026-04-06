package mbjson

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerJSONCommands(m *Module, r runtime.Registrar) {
	r.Register("JSON.PARSE", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jParse(m, rt, args...) })
	r.Register("JSON.LOADFILE", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jParse(m, rt, args...) })
	r.Register("JSON.PARSESTRING", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jParseString(m, rt, args...) })
	r.Register("JSON.MAKE", "json", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return jMake(m, a) }))
	r.Register("JSON.MAKEARRAY", "json", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return jMakeArray(m, a) }))
	r.Register("JSON.FREE", "json", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return jFree(m, a) }))

	r.Register("JSON.HAS", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jHas(m, rt, args...) })
	r.Register("JSON.TYPE$", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jTypeStr(m, rt, args...) })
	r.Register("JSON.LEN", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jLen(m, rt, args...) })
	r.Register("JSON.KEYS", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jKeys(m, rt, args...) })

	r.Register("JSON.GETSTRING", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetString(m, rt, args...) })
	r.Register("JSON.GETINT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetInt(m, rt, args...) })
	r.Register("JSON.GETFLOAT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetFloat(m, rt, args...) })
	r.Register("JSON.GETBOOL", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetBool(m, rt, args...) })
	r.Register("JSON.GETARRAY", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetArray(m, rt, args...) })
	r.Register("JSON.GETOBJECT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jGetObject(m, rt, args...) })

	r.Register("JSON.SETSTRING", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetString(m, rt, args...) })
	r.Register("JSON.SETINT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetInt(m, rt, args...) })
	r.Register("JSON.SETFLOAT", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetFloat(m, rt, args...) })
	r.Register("JSON.SETBOOL", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetBool(m, rt, args...) })
	r.Register("JSON.SETNULL", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jSetNull(m, rt, args...) })
	r.Register("JSON.DELETE", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jDelete(m, rt, args...) })
	r.Register("JSON.CLEAR", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jClear(m, rt, args...) })
	r.Register("JSON.APPEND", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jAppend(m, rt, args...) })

	r.Register("JSON.TOSTRING", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jToString(m, rt, args...) })
	r.Register("JSON.PRETTY", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jPretty(m, rt, args...) })
	r.Register("JSON.MINIFY", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jMinify(m, rt, args...) })
	r.Register("JSON.TOFILE", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jToFile(m, rt, args...) })
	r.Register("JSON.SAVEFILE", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jToFile(m, rt, args...) })
	r.Register("JSON.TOFILEPRETTY", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jToFilePretty(m, rt, args...) })
	r.Register("JSON.TOCSV", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jToCSV(m, rt, args...) })

	r.Register("JSON.QUERY", "json", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return jQuery(m, rt, args...) })
}
