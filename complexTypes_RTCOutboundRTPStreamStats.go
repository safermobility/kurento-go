package kurento

type RTCOutboundRTPStreamStats struct {
	PacketsSent   int64
	BytesSent     int64
	TargetBitrate float64
	RoundTripTime float64
}
