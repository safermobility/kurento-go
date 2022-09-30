package kurento

// How to fix gaps when they are found in the recorded stream.
// <p>
// Gaps are typically caused by packet loss in the input streams, such as when an
// RTP or WebRTC media flow suffers from network congestion and some packets don't
// arrive at the media server.
// </p>
// <p>Different ways of handling gaps have different tradeoffs:</p>
// <ul>
// <li>
// <strong>NONE</strong>: Do not fix gaps.
// <p>
// Leave the stream as-is, and store it with any gaps that the stream might
// have. Some players are clever enough to adapt to this during playback, so
// that the gaps are reduced to a minimum and no problems are perceived by
// the user; other players are not so sophisticated, and will struggle trying
// to decode a file that contains gaps. For example, trying to play such a
// file directly with Chrome will cause lipsync issues (audio and video will
// fall out of sync).
// </p>
// <p>
// This is the best choice if you need consistent durations across multiple
// simultaneous recordings, or if you are anyway going to post-process the
// recordings (e.g. with an extra FFmpeg step).
// </p>
// <p>
// For example, assume a session length of 15 seconds: packets arrive
// correctly during the first 5 seconds, then there is a gap, then data
// arrives again for the last 5 seconds. Also, for simplicity, assume 1 frame
// per second. With no fix for gaps, the RecorderEndpoint will store each
// frame as-is, with these timestamps:
// </p>
// <pre>
// frame 1  - 00:01
// frame 2  - 00:02
// frame 3  - 00:03
// frame 4  - 00:04
// frame 5  - 00:05
// frame 11 - 00:11
// frame 12 - 00:12
// frame 13 - 00:13
// frame 14 - 00:14
// frame 15 - 00:15
// </pre>
// <p>
// Notice how the frames between 6 to 10 are missing, but the last 5 frames
// still conserve their original timestamp. The total length of the file is
// detected as 15 seconds by most players, although playback could stutter or
// hang during the missing section.
// </p>
// </li>
// <li>
// <strong>GENPTS</strong>: Adjust timestamps to generate a smooth progression
// over all frames.
// <p>
// This technique rewrites the timestamp of all frames, so that gaps are
// suppressed. It provides the best playback experience for recordings that
// need to be played as-is (i.e. they won't be post-processed). However,
// fixing timestamps might cause a change in the total duration of a file. So
// different recordings from the same session might end up with slightly
// different durations.
// </p>
// <p>
// In our example, the RecorderEndpoint will change all timestamps that
// follow a gap in the stream, and store each frame as follows:
// </p>
// <pre>
// frame 1  - 00:01
// frame 2  - 00:02
// frame 3  - 00:03
// frame 4  - 00:04
// frame 5  - 00:05
// frame 11 - 00:06
// frame 12 - 00:07
// frame 13 - 00:08
// frame 14 - 00:09
// frame 15 - 00:10
// </pre>
// <p>
// Notice how the frames between 6 to 10 are missing, and the last 5 frames
// have their timestamps corrected to provide a smooth increment over the
// previous ones. The total length of the file is detected as 10 seconds, and
// playback should be correct throughout the whole file.
// </p>
// </li>
// <li>
// <strong>FILL_IF_TRANSCODING</strong>: (NOT IMPLEMENTED YET).
// <p>This is a proposal for future improvement of the RecorderEndpoint.</p>
// <p>
// It is possible to perform a dynamic adaptation of audio rate and add frame
// duplication to the video, such that the missing parts are filled with
// artificial data. This has the advantage of providing a smooth playback
// result, and at the same time conserving all original timestamps.
// </p>
// <p>
// However, the main issue with this method is that it requires accessing the
// decoded media; i.e., transcoding must be active. For this reason, the
// proposal is to offer this option to be enabled only when transcoding would
// still happen anyways.
// </p>
// <p>
// In our example, the RecorderEndpoint would change all missing frames like
// this:
// </p>
// <pre>
// frame 1  - 00:01
// frame 2  - 00:02
// frame 3  - 00:03
// frame 4  - 00:04
// frame 5  - 00:05
// fake frame - 00:06
// fake frame - 00:07
// fake frame - 00:08
// fake frame - 00:09
// fake frame - 00:10
// frame 11 - 00:11
// frame 12 - 00:12
// frame 13 - 00:13
// frame 14 - 00:14
// frame 15 - 00:15
// </pre>
// <p>
// This joins the best of both worlds: on one hand, the playback should be
// smooth and even the most basic players should be able to handle the
// recording files without issue. On the other, the total length of the file
// is left unmodified, so it matches with the expected duration of the
// sessions that are being recorded.
// </p>
// </li>
// </ul>
//
type GapsFixMethod string

// Implement fmt.Stringer interface
func (t GapsFixMethod) String() string {
	return string(t)
}

const (
	GAPSFIXMETHOD_NONE                GapsFixMethod = "NONE"
	GAPSFIXMETHOD_GENPTS              GapsFixMethod = "GENPTS"
	GAPSFIXMETHOD_FILL_IF_TRANSCODING GapsFixMethod = "FILL_IF_TRANSCODING"
)