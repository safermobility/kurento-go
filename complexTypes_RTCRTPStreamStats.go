package kurento

type RTCRTPStreamStats struct {
	Ssrc             string
	AssociateStatsId string
	IsRemote         bool
	MediaTrackId     string
	TransportId      string
	CodecId          string
	FirCount         int64
	PliCount         int64
	NackCount        int64
	SliCount         int64
	Remb             int64
	PacketsLost      int64
	FractionLost     float64
}
