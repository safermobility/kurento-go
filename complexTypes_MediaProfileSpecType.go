package kurento

// Media profile, used by the RecorderEndpoint builder to specify the codecs and media container that should be used for the recordings.
type MediaProfileSpecType string

// Implement fmt.Stringer interface
func (t MediaProfileSpecType) String() string {
	return string(t)
}

const (
	MEDIAPROFILESPECTYPE_WEBM                   MediaProfileSpecType = "WEBM"
	MEDIAPROFILESPECTYPE_MKV                    MediaProfileSpecType = "MKV"
	MEDIAPROFILESPECTYPE_MP4                    MediaProfileSpecType = "MP4"
	MEDIAPROFILESPECTYPE_WEBM_VIDEO_ONLY        MediaProfileSpecType = "WEBM_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_WEBM_AUDIO_ONLY        MediaProfileSpecType = "WEBM_AUDIO_ONLY"
	MEDIAPROFILESPECTYPE_MKV_VIDEO_ONLY         MediaProfileSpecType = "MKV_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_MKV_AUDIO_ONLY         MediaProfileSpecType = "MKV_AUDIO_ONLY"
	MEDIAPROFILESPECTYPE_MP4_VIDEO_ONLY         MediaProfileSpecType = "MP4_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_MP4_AUDIO_ONLY         MediaProfileSpecType = "MP4_AUDIO_ONLY"
	MEDIAPROFILESPECTYPE_JPEG_VIDEO_ONLY        MediaProfileSpecType = "JPEG_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_KURENTO_SPLIT_RECORDER MediaProfileSpecType = "KURENTO_SPLIT_RECORDER"
	MEDIAPROFILESPECTYPE_FLV                    MediaProfileSpecType = "FLV"
)
