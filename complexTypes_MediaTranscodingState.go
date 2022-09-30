package kurento

// Transcoding state for a media.
type MediaTranscodingState string

// Implement fmt.Stringer interface
func (t MediaTranscodingState) String() string {
	return string(t)
}

const (
	MEDIATRANSCODINGSTATE_TRANSCODING     MediaTranscodingState = "TRANSCODING"
	MEDIATRANSCODINGSTATE_NOT_TRANSCODING MediaTranscodingState = "NOT_TRANSCODING"
)
