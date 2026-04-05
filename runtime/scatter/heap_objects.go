package scatter

import "moonbasic/vm/heap"

type ScatterObject struct {
	Name string
	X    []float32
	Y    []float32
	Z    []float32
	Seed int64
	freed bool
}

func (s *ScatterObject) TypeName() string { return "Scatter" }
func (s *ScatterObject) TypeTag() uint16   { return heap.TagScatterSet }
func (s *ScatterObject) Free()             { s.freed = true; s.X, s.Y, s.Z = nil, nil, nil }

type PropObject struct {
	X, Y, Z float32
	freed   bool
}

func (p *PropObject) TypeName() string { return "Prop" }
func (p *PropObject) TypeTag() uint16  { return heap.TagProp }
func (p *PropObject) Free()            { p.freed = true }
