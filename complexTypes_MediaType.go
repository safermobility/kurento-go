// Code generated by kurento-go-generator. DO NOT EDIT.

package kurento

// Type of media stream to be exchanged.
// Can take the values AUDIO, DATA or VIDEO.
type MediaType string

// Implement fmt.Stringer interface
func (t MediaType) String() string {
	return string(t)
}

const (
	MEDIATYPE_AUDIO MediaType = "AUDIO"
	MEDIATYPE_DATA  MediaType = "DATA"
	MEDIATYPE_VIDEO MediaType = "VIDEO"
)
