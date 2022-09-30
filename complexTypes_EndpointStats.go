package kurento

type EndpointStats struct {
	AudioE2ELatency float64
	VideoE2ELatency float64
	E2ELatency      []MediaLatencyStat
}
