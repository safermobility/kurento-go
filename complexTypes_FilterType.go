package kurento

// Type of filter to be created.
// Can take the values AUDIO, VIDEO or AUTODETECT.
type FilterType string

// Implement fmt.Stringer interface
func (t FilterType) String() string {
	return string(t)
}

const (
	FILTERTYPE_AUDIO      FilterType = "AUDIO"
	FILTERTYPE_AUTODETECT FilterType = "AUTODETECT"
	FILTERTYPE_VIDEO      FilterType = "VIDEO"
)
