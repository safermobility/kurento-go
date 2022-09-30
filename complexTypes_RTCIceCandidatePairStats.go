package kurento

type RTCIceCandidatePairStats struct {
	TransportId              string
	LocalCandidateId         string
	RemoteCandidateId        string
	State                    RTCStatsIceCandidatePairState
	Priority                 int64
	Nominated                bool
	Writable                 bool
	Readable                 bool
	BytesSent                int64
	BytesReceived            int64
	RoundTripTime            float64
	AvailableOutgoingBitrate float64
	AvailableIncomingBitrate float64
}
