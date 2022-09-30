package kurento

type RTCCodec struct {
	PayloadType int64
	Codec       string
	ClockRate   int64
	Channels    int64
	Parameters  string
}
