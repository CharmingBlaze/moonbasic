//go:build !cgo

package weathermod

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerWeather(m *Module, r runtime.Registrar) {
	h := func(n string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO", n)
		}
	}
	r.Register("WEATHER.MAKE", "weather", h("WEATHER.MAKE"))
	r.Register("WEATHER.FREE", "weather", h("WEATHER.FREE"))
	r.Register("WEATHER.UPDATE", "weather", h("WEATHER.UPDATE"))
	r.Register("WEATHER.DRAW", "weather", h("WEATHER.DRAW"))
	r.Register("WEATHER.SETTYPE", "weather", h("WEATHER.SETTYPE"))
	r.Register("WEATHER.GETCOVERAGE", "weather", h("WEATHER.GETCOVERAGE"))
	r.Register("WEATHER.GETTYPE", "weather", h("WEATHER.GETTYPE"))
	r.Register("FOG.ENABLE", "fog", h("FOG.ENABLE"))
	r.Register("FOG.SETNEAR", "fog", h("FOG.SETNEAR"))
	r.Register("FOG.SETFAR", "fog", h("FOG.SETFAR"))
	r.Register("FOG.SETCOLOR", "fog", h("FOG.SETCOLOR"))
	r.Register("WIND.SET", "wind", h("WIND.SET"))
	r.Register("WIND.GETSTRENGTH", "wind", h("WIND.GETSTRENGTH"))
}
