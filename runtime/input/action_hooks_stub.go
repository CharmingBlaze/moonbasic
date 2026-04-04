//go:build !cgo

package input

type actionQuery struct{}

func (actionQuery) keyPressed(k int32) bool         { _ = k; return false }
func (actionQuery) keyDown(k int32) bool            { _ = k; return false }
func (actionQuery) keyReleased(k int32) bool       { _ = k; return false }
func (actionQuery) gamepadBtnPressed(pad, btn int32) bool {
	_, _ = pad, btn
	return false
}
func (actionQuery) gamepadBtnDown(pad, btn int32) bool {
	_, _ = pad, btn
	return false
}
func (actionQuery) gamepadBtnReleased(pad, btn int32) bool {
	_, _ = pad, btn
	return false
}
func (actionQuery) gamepadAxis(pad, axis int32) float32 {
	_, _ = pad, axis
	return 0
}

func actionQueries() actionQuery { return actionQuery{} }
