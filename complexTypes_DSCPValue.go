package kurento

// Possible DSCP values
// <p>
// WebRTC recommended values are taken from RFC 8837 https://datatracker.ietf.org/doc/html/rfc8837#section-5 , These are the values from AUDIO_VERYLOW to DATA_HIGH. First element in the
// name indicates kind of traffic that it should apply to, the second indicates relative priority. For video, a third field would indicate if the traffic is intended for high throughput or not.
// As indicated on RFC 8837 section 5 diagram:
//
// +=======================+==========+=====+============+============+
// |       Flow Type       | Very Low | Low |   Medium   |    High    |
// +=======================+==========+=====+============+============+
// |         Audio         |  LE (1)  |  DF |  EF (46)   |  EF (46)   |
// |                       |          | (0) |            |            |
// +-----------------------+----------+-----+------------+------------+
// +-----------------------+----------+-----+------------+------------+
// |   Interactive Video   |  LE (1)  |  DF | AF42, AF43 | AF41, AF42 |
// | with or without Audio |          | (0) |  (36, 38)  |  (34, 36)  |
// +-----------------------+----------+-----+------------+------------+
// +-----------------------+----------+-----+------------+------------+
// | Non-Interactive Video |  LE (1)  |  DF | AF32, AF33 | AF31, AF32 |
// | with or without Audio |          | (0) |  (28, 30)  |  (26, 28)  |
// +-----------------------+----------+-----+------------+------------+
// +-----------------------+----------+-----+------------+------------+
// |          Data         |  LE (1)  |  DF |    AF11    |    AF21    |
// |                       |          | (0) |            |            |
// +-----------------------+----------+-----+------------+------------+
//
// As indicated also in RFC, non interactive video is not considered
// <p>
// Apart from the WebRTC recommended values, we also include all possible values are referenced in http://www.iana.org/assignments/dscp-registry/dscp-registry.xml of course some of
// those values are synonyms for the WebRTC recommended ones, they are included mainly for completeness
// <p>
// And as a last point, we include a shorthand for Chrome supported markings for low  (CS0), very low (CS1), medium (CS7) and high (CS7) priorities in priority property for RTCRtpSender parameters. See https://developer.mozilla.org/en-US/docs/Web/API/RTCRtpSender/setParameters
// <p>
// This only covers outgoing network packets from KMS, to complete the solution, DSCP must be also requested on client, unfortunately for traffic on the other direction, this must be requested to the
// browser or client. On browser, the client application needs to use the following API https://www.w3.org/TR/webrtc-priority/
type DSCPValue string

// Implement fmt.Stringer interface
func (t DSCPValue) String() string {
	return string(t)
}

const (
	DSCPVALUE_NO_DSCP                 DSCPValue = "NO_DSCP"
	DSCPVALUE_NO_VALUE                DSCPValue = "NO_VALUE"
	DSCPVALUE_AUDIO_VERYLOW           DSCPValue = "AUDIO_VERYLOW"
	DSCPVALUE_AUDIO_LOW               DSCPValue = "AUDIO_LOW"
	DSCPVALUE_AUDIO_MEDIUM            DSCPValue = "AUDIO_MEDIUM"
	DSCPVALUE_AUDIO_HIGH              DSCPValue = "AUDIO_HIGH"
	DSCPVALUE_VIDEO_VERYLOW           DSCPValue = "VIDEO_VERYLOW"
	DSCPVALUE_VIDEO_LOW               DSCPValue = "VIDEO_LOW"
	DSCPVALUE_VIDEO_MEDIUM            DSCPValue = "VIDEO_MEDIUM"
	DSCPVALUE_VIDEO_MEDIUM_THROUGHPUT DSCPValue = "VIDEO_MEDIUM_THROUGHPUT"
	DSCPVALUE_VIDEO_HIGH              DSCPValue = "VIDEO_HIGH"
	DSCPVALUE_VIDEO_HIGH_THROUGHPUT   DSCPValue = "VIDEO_HIGH_THROUGHPUT"
	DSCPVALUE_DATA_VERYLOW            DSCPValue = "DATA_VERYLOW"
	DSCPVALUE_DATA_LOW                DSCPValue = "DATA_LOW"
	DSCPVALUE_DATA_MEDIUM             DSCPValue = "DATA_MEDIUM"
	DSCPVALUE_DATA_HIGH               DSCPValue = "DATA_HIGH"
	DSCPVALUE_CHROME_HIGH             DSCPValue = "CHROME_HIGH"
	DSCPVALUE_CHROME_MEDIUM           DSCPValue = "CHROME_MEDIUM"
	DSCPVALUE_CHROME_LOW              DSCPValue = "CHROME_LOW"
	DSCPVALUE_CHROME_VERYLOW          DSCPValue = "CHROME_VERYLOW"
	DSCPVALUE_CS0                     DSCPValue = "CS0"
	DSCPVALUE_CS1                     DSCPValue = "CS1"
	DSCPVALUE_CS2                     DSCPValue = "CS2"
	DSCPVALUE_CS3                     DSCPValue = "CS3"
	DSCPVALUE_CS4                     DSCPValue = "CS4"
	DSCPVALUE_CS5                     DSCPValue = "CS5"
	DSCPVALUE_CS6                     DSCPValue = "CS6"
	DSCPVALUE_CS7                     DSCPValue = "CS7"
	DSCPVALUE_AF11                    DSCPValue = "AF11"
	DSCPVALUE_AF12                    DSCPValue = "AF12"
	DSCPVALUE_AF13                    DSCPValue = "AF13"
	DSCPVALUE_AF21                    DSCPValue = "AF21"
	DSCPVALUE_AF22                    DSCPValue = "AF22"
	DSCPVALUE_AF23                    DSCPValue = "AF23"
	DSCPVALUE_AF31                    DSCPValue = "AF31"
	DSCPVALUE_AF32                    DSCPValue = "AF32"
	DSCPVALUE_AF33                    DSCPValue = "AF33"
	DSCPVALUE_AF41                    DSCPValue = "AF41"
	DSCPVALUE_AF42                    DSCPValue = "AF42"
	DSCPVALUE_AF43                    DSCPValue = "AF43"
	DSCPVALUE_EF                      DSCPValue = "EF"
	DSCPVALUE_VOICEADMIT              DSCPValue = "VOICEADMIT"
	DSCPVALUE_LE                      DSCPValue = "LE"
)