package cpu2a03

// InputDvice is a generic device connected to either port of the NES
// Typically a controller, but can be other stuff
type InputDevice interface {
	Read() byte // CLK trigger + read
	Write(byte) // OUT0, OUT1 and OUT2 update
}
