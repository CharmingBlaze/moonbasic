package weathermod

import (
	"moonbasic/runtime"
)

func (m *Module) Register(r runtime.Registrar) {
	r.Register("WEATHER.MAKE", "weather", runtime.AdaptLegacy(m.wMake))
	r.Register("WEATHER.FREE", "weather", runtime.AdaptLegacy(m.wFree))
	r.Register("WEATHER.UPDATE", "weather", runtime.AdaptLegacy(m.wUpdate))
	r.Register("WEATHER.DRAW", "weather", runtime.AdaptLegacy(m.wDraw))
	r.Register("WEATHER.SETTYPE", "weather", runtime.AdaptLegacy(m.wSetType))
	r.Register("WEATHER.GETCOVERAGE", "weather", runtime.AdaptLegacy(m.wGetCoverage))
	r.Register("WEATHER.GETTYPE", "weather", runtime.AdaptLegacy(m.wGetType))

	r.Register("FOG.ENABLE", "fog", runtime.AdaptLegacy(m.fogEnable))
	r.Register("FOG.SETNEAR", "fog", runtime.AdaptLegacy(m.fogSetNear))
	r.Register("FOG.SETFAR", "fog", runtime.AdaptLegacy(m.fogSetFar))
	r.Register("FOG.SETCOLOR", "fog", runtime.AdaptLegacy(m.fogSetColor))
	r.Register("FOG.SETRANGE", "fog", runtime.AdaptLegacy(m.fogSetRange))

	r.Register("WIND.SET", "wind", runtime.AdaptLegacy(m.windSet))
	r.Register("WIND.GETSTRENGTH", "wind", runtime.AdaptLegacy(m.windGetStrength))
}
