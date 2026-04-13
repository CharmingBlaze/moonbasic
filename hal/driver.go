package hal

// Driver aggregates all hardware-dependent subsystems.
type Driver struct {
	Video  VideoDevice
	Input  InputDevice
	System SystemDevice
}
