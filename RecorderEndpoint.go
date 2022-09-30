package kurento

import "fmt"

type IRecorderEndpoint interface {
	Record() error
	StopAndWait() error
}

// Provides functionality to store media contents.
// <p>
// RecorderEndpoint can store media into local files or send it to a remote
// network storage. When another `MediaElement` is connected to a
// RecorderEndpoint, the media coming from the former will be muxed into
// the selected recording format and stored in the designated location.
// </p>
// <p>
// These parameters must be provided to create a RecorderEndpoint, and they
// cannot be changed afterwards:
// </p>
// <ul>
// <li>
// <strong>Destination URI</strong>, where media will be stored. These formats
// are supported:
// <ul>
// <li>
// File: A file path that will be written into the local file system.
// Example:
// <ul>
// <li><code>file:///path/to/file</code></li>
// </ul>
// </li>
// <li>
// HTTP: A POST request will be used against a remote server. The server
// must support using the <i>chunked</i> encoding mode (HTTP header
// <code>Transfer-Encoding: chunked</code>). Examples:
// <ul>
// <li><code>http(s)://{server-ip}/path/to/file</code></li>
// <li>
// <code>
// http(s)://{username}:{password}@{server-ip}:{server-port}/path/to/file
// </code>
// </li>
// </ul>
// </li>
// <li>
// Relative URIs (with no schema) are supported. They are completed by
// prepending a default URI defined by property <i>defaultPath</i>. This
// property is defined in the configuration file
// <i>/etc/kurento/modules/kurento/UriEndpoint.conf.ini</i>, and the
// default value is <code>file:///var/lib/kurento/</code>
// </li>
// <li>
// <strong>
// NOTE (for current versions of Kurento 6.x): special characters are not
// supported in <code>{username}</code> or <code>{password}</code>.
// </strong>
// This means that <code>{username}</code> cannot contain colons
// (<code>:</code>), and <code>{password}</code> cannot contain 'at' signs
// (<code>@</code>). This is a limitation of GStreamer 1.8 (the underlying
// media framework behind Kurento), and is already fixed in newer versions
// (which the upcoming Kurento 7.x will use).
// </li>
// <li>
// <strong>
// NOTE (for upcoming Kurento 7.x): special characters in
// <code>{username}</code> or <code>{password}</code> must be
// url-encoded.
// </strong>
// This means that colons (<code>:</code>) should be replaced with
// '<code>%3A</code>', and 'at' signs (<code>@</code>) should be replaced
// with '<code>%40</code>'.
// </li>
// </ul>
// </li>
// <li>
// <strong>Media Profile</strong> (:rom:enum:`MediaProfileSpecType`), which
// determines the video and audio encoding. See below for more details.
// </li>
// <li>
// <strong>EndOfStream</strong> (optional), a parameter that dictates if the
// recording should be automatically stopped once the EOS event is detected.
// </li>
// </ul>
// <p>
// Note that
// <strong>
// RecorderEndpoint requires write permissions to the destination
// </strong>
// ; otherwise, the media server won't be able to store any information, and an
// :rom:evt:`Error` will be fired. Make sure your application subscribes to this
// event, otherwise troubleshooting issues will be difficult.
// </p>
// <ul>
// <li>
// To write local files (if you use <code>file://</code>), the system user that
// is owner of the media server process needs to have write permissions for the
// requested path. By default, this user is named '<code>kurento</code>'.
// </li>
// <li>
// To record through HTTP, the remote server must be accessible through the
// network, and also have the correct write permissions for the destination
// path.
// </li>
// </ul>
// <p>
// Recording will start as soon as the user invokes the
// <code>record()</code> method. The recorder will then store, in the location
// indicated, the media that the source is sending to the endpoint. If no media
// is being received, or no endpoint has been connected, then the destination
// will be empty. The recorder starts storing information into the file as soon
// as it gets it.
// </p>
// <p>
// <strong>Recording must be stopped</strong> when no more data should be stored.
// This is done with the <code>stopAndWait()</code> method, which blocks and
// returns only after all the information was stored correctly.
// </p>
// <p>
// The source endpoint can be hot-swapped while the recording is taking place.
// The recorded file will then contain different feeds. When switching video
// sources, if the new video has different size, the recorder will retain the
// size of the previous source. If the source is disconnected, the last frame
// recorded will be shown for the duration of the disconnection, or until the
// recording is stopped.
// </p>
// <p>
// <strong>
// NOTE: It is recommended to start recording only after media arrives.
// </strong>
// For this, you may use the <code>MediaFlowInStateChanged</code> and
// <code>MediaFlowOutStateChanged</code>
// events of your endpoints, and synchronize the recording with the moment media
// comes into the Recorder.
// </p>
// <p>
// <strong>
// WARNING: All connected media types must be flowing to the RecorderEndpoint.
// </strong>
// If you used the default <code>connect()</code> method, it will assume both
// AUDIO and VIDEO. Failing to provide both kinds of media will result in the
// RecorderEndpoint creating an empty file and buffering indefinitely; the
// recorder waits until all kinds of media start arriving, in order to
// synchronize them appropriately.<br>
// For audio-only or video-only recordings, make sure to use the correct,
// media-specific variant of the <code>connect()</code> method.
// </p>
// <p>
// For example:
// </p>
// <ol>
// <li>
// When a web browser's video arrives to Kurento via WebRTC, your
// WebRtcEndpoint will emit a <code>MediaFlowOutStateChanged</code> event.
// </li>
// <li>
// When video starts flowing from the WebRtcEndpoint to the RecorderEndpoint,
// the RecorderEndpoint will emit a <code>MediaFlowInStateChanged</code> event.
// You should start recording at this point.
// </li>
// <li>
// You should only start recording when RecorderEndpoint has notified a
// <code>MediaFlowInStateChanged</code> for ALL streams. So, if you record
// AUDIO+VIDEO, your application must receive a
// <code>MediaFlowInStateChanged</code> event for audio, and another
// <code>MediaFlowInStateChanged</code> event for video.
// </li>
// </ol>
//
type RecorderEndpoint struct {
	UriEndpoint
}

// Return contructor params to be called by "Create".
func (elem *RecorderEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":     fmt.Sprintf("%s", from),
		"uri":               "",
		"mediaProfile":      fmt.Sprintf("%s", from),
		"stopOnEndOfStream": false,
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Starts storing media received through the sink pad.
func (elem *RecorderEndpoint) Record() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "record",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil {
		return fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}
	return nil

}

// Stops recording and does not return until all the content has been written to the selected uri. This can cause timeouts on some clients if there is too much content to write, or the transport is slow
func (elem *RecorderEndpoint) StopAndWait() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "stopAndWait",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil {
		return fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}
	return nil

}