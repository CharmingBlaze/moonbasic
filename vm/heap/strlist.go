package heap

// StringList is a heap-owned slice of strings produced by SPLIT$ and consumed by JOIN$.
type StringList struct {
	Items []string
}

func (s *StringList) TypeName() string { return "StringList" }

func (s *StringList) TypeTag() uint16 { return TagStringList }

func (s *StringList) Free() {
	s.Items = nil
}
