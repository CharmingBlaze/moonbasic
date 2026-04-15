package semantic

import "fmt"

// DeprecationNotice records one use of a deprecated built-in alias at a source location.
type DeprecationNotice struct {
	File           string
	Line, Col      int
	DeprecatedKey  string
	ReplacementKey string
}

// String matches the CLI stderr warning format.
func (n DeprecationNotice) String() string {
	return fmt.Sprintf("%s:%d:%d: deprecated command %s; use %s",
		n.File, n.Line, n.Col, n.DeprecatedKey, n.ReplacementKey)
}
