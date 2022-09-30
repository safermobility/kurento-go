package kurento

// States of an ICE component.
type IceComponentState string

// Implement fmt.Stringer interface
func (t IceComponentState) String() string {
	return string(t)
}

const (
	ICECOMPONENTSTATE_DISCONNECTED IceComponentState = "DISCONNECTED"
	ICECOMPONENTSTATE_GATHERING    IceComponentState = "GATHERING"
	ICECOMPONENTSTATE_CONNECTING   IceComponentState = "CONNECTING"
	ICECOMPONENTSTATE_CONNECTED    IceComponentState = "CONNECTED"
	ICECOMPONENTSTATE_READY        IceComponentState = "READY"
	ICECOMPONENTSTATE_FAILED       IceComponentState = "FAILED"
)
