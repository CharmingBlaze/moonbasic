package sky

import "moonbasic/vm/heap"

// SkyObject holds day/night cycle state (drawn as large sphere each frame).
type SkyObject struct {
	Time      float32
	DayLength float32
	freed     bool
}

func (s *SkyObject) TypeName() string { return "Sky" }
func (s *SkyObject) TypeTag() uint16  { return heap.TagSky }

func (s *SkyObject) Free() {
	s.freed = true
}
