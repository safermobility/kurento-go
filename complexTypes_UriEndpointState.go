package kurento

// State of the endpoint
type UriEndpointState string

// Implement fmt.Stringer interface
func (t UriEndpointState) String() string {
	return string(t)
}

const (
	URIENDPOINTSTATE_STOP  UriEndpointState = "STOP"
	URIENDPOINTSTATE_START UriEndpointState = "START"
	URIENDPOINTSTATE_PAUSE UriEndpointState = "PAUSE"
)
