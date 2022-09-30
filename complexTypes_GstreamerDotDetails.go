package kurento

// Details of gstreamer dot graphs
type GstreamerDotDetails string

// Implement fmt.Stringer interface
func (t GstreamerDotDetails) String() string {
	return string(t)
}

const (
	GSTREAMERDOTDETAILS_SHOW_MEDIA_TYPE         GstreamerDotDetails = "SHOW_MEDIA_TYPE"
	GSTREAMERDOTDETAILS_SHOW_CAPS_DETAILS       GstreamerDotDetails = "SHOW_CAPS_DETAILS"
	GSTREAMERDOTDETAILS_SHOW_NON_DEFAULT_PARAMS GstreamerDotDetails = "SHOW_NON_DEFAULT_PARAMS"
	GSTREAMERDOTDETAILS_SHOW_STATES             GstreamerDotDetails = "SHOW_STATES"
	GSTREAMERDOTDETAILS_SHOW_FULL_PARAMS        GstreamerDotDetails = "SHOW_FULL_PARAMS"
	GSTREAMERDOTDETAILS_SHOW_ALL                GstreamerDotDetails = "SHOW_ALL"
	GSTREAMERDOTDETAILS_SHOW_VERBOSE            GstreamerDotDetails = "SHOW_VERBOSE"
)
