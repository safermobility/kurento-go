package kurento

type RTCTransportStats struct {
	BytesSent               int64
	BytesReceived           int64
	RtcpTransportStatsId    string
	ActiveConnection        bool
	SelectedCandidatePairId string
	LocalCertificateId      string
	RemoteCertificateId     string
}
