package kurento

type RTCDataChannelStats struct {
	Label            string
	Protocol         string
	Datachannelid    int64
	State            RTCDataChannelState
	MessagesSent     int64
	BytesSent        int64
	MessagesReceived int64
	BytesReceived    int64
}
