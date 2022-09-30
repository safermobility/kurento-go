package kurento

type IceCandidate struct {
	Candidate     string
	SdpMid        string
	SdpMLineIndex int
}

func (t IceCandidate) CustomSerialize() map[string]interface{} {
	ret := make(map[string]interface{})

	ret["candidate"] = t.Candidate

	ret["sdpMid"] = t.SdpMid

	ret["sdpMLineIndex"] = t.SdpMLineIndex

	ret["__type__"] = "IceCandidate"
	ret["__module__"] = "kurento"
	return ret
}
