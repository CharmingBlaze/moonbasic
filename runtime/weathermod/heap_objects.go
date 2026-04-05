package weathermod

import "moonbasic/vm/heap"

type WeatherObject struct {
	Kind      string
	Intensity float32
	Coverage  float32
	freed     bool
}

func (w *WeatherObject) TypeName() string { return "Weather" }
func (w *WeatherObject) TypeTag() uint16   { return heap.TagWeather }
func (w *WeatherObject) Free()             { w.freed = true }
