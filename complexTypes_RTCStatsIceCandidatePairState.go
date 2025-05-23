// Code generated by kurento-go-generator. DO NOT EDIT.

package kurento

// Represents the state of the checklist for the local and remote candidates in a pair.
type RTCStatsIceCandidatePairState string

// Implement fmt.Stringer interface
func (t RTCStatsIceCandidatePairState) String() string {
	return string(t)
}

const (
	RTCSTATSICECANDIDATEPAIRSTATE_frozen     RTCStatsIceCandidatePairState = "frozen"
	RTCSTATSICECANDIDATEPAIRSTATE_waiting    RTCStatsIceCandidatePairState = "waiting"
	RTCSTATSICECANDIDATEPAIRSTATE_inprogress RTCStatsIceCandidatePairState = "inprogress"
	RTCSTATSICECANDIDATEPAIRSTATE_failed     RTCStatsIceCandidatePairState = "failed"
	RTCSTATSICECANDIDATEPAIRSTATE_succeeded  RTCStatsIceCandidatePairState = "succeeded"
	RTCSTATSICECANDIDATEPAIRSTATE_cancelled  RTCStatsIceCandidatePairState = "cancelled"
)
