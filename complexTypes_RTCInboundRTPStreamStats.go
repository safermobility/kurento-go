package kurento

type RTCInboundRTPStreamStats struct {
	PacketsReceived int64
	BytesReceived   int64
	Jitter          float64
}
