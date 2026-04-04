package heap

// ArrayFlatLen returns the total number of elements in a heap array handle, or -1 if invalid.
func (s *Store) ArrayFlatLen(h Handle) int {
	a, err := Cast[*Array](s, h)
	if err != nil {
		return -1
	}
	return a.totalLen()
}

// ArrayGetFloat returns one numeric cell from a heap array using a flat index.
func (s *Store) ArrayGetFloat(h Handle, flat int64) (float64, bool) {
	a, err := Cast[*Array](s, h)
	if err != nil {
		return 0, false
	}
	v, err := a.Get([]int64{flat})
	return v, err == nil
}
