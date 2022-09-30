package kurento

// State of the media.
type MediaState string

// Implement fmt.Stringer interface
func (t MediaState) String() string {
	return string(t)
}

const (
	MEDIASTATE_DISCONNECTED MediaState = "DISCONNECTED"
	MEDIASTATE_CONNECTED    MediaState = "CONNECTED"
)
