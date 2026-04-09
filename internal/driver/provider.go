package driver

// WindowProvider is the future hook for consolidating Raylib backends (CGO vs purego sidecar).
// Today, [GetDefaultDriver] plus [moonbasic/runtime/window.Module.BindDriverSelection] perform selection.
type WindowProvider interface {
	DriverKind() Kind
}
