package biome

import "moonbasic/vm/heap"

type BiomeObject struct {
	Name      string
	TempC     float32
	Humidity  float32
	freed     bool
}

func (b *BiomeObject) TypeName() string { return "Biome" }
func (b *BiomeObject) TypeTag() uint16  { return heap.TagBiome }
func (b *BiomeObject) Free()            { b.freed = true }
