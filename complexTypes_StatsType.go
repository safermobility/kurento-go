package kurento

// The type of the object.
type StatsType string

// Implement fmt.Stringer interface
func (t StatsType) String() string {
	return string(t)
}

const (
	STATSTYPE_inboundrtp      StatsType = "inboundrtp"
	STATSTYPE_outboundrtp     StatsType = "outboundrtp"
	STATSTYPE_session         StatsType = "session"
	STATSTYPE_datachannel     StatsType = "datachannel"
	STATSTYPE_track           StatsType = "track"
	STATSTYPE_transport       StatsType = "transport"
	STATSTYPE_candidatepair   StatsType = "candidatepair"
	STATSTYPE_localcandidate  StatsType = "localcandidate"
	STATSTYPE_remotecandidate StatsType = "remotecandidate"
	STATSTYPE_element         StatsType = "element"
	STATSTYPE_endpoint        StatsType = "endpoint"
)
