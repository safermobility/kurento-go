package kurento

// Codec used for transmission of audio.
type AudioCodec string

// Implement fmt.Stringer interface
func (t AudioCodec) String() string {
	return string(t)
}

const (
	AUDIOCODEC_OPUS AudioCodec = "OPUS"
	AUDIOCODEC_PCMU AudioCodec = "PCMU"
	AUDIOCODEC_RAW  AudioCodec = "RAW"
)
