package kurento

type ElementStats struct {
	InputAudioLatency float64
	InputVideoLatency float64
	InputLatency      []MediaLatencyStat
}
