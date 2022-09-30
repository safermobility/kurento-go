package kurento

type RTCMediaStreamTrackStats struct {
	TrackIdentifier           string
	RemoteSource              bool
	SsrcIds                   []string
	FrameWidth                int64
	FrameHeight               int64
	FramesPerSecond           float64
	FramesSent                int64
	FramesReceived            int64
	FramesDecoded             int64
	FramesDropped             int64
	FramesCorrupted           int64
	AudioLevel                float64
	EchoReturnLoss            float64
	EchoReturnLossEnhancement float64
}
