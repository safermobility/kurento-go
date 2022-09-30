package kurento

import "fmt"

type IPlayerEndpoint interface {
	Play() error
}

// Retrieves content from external sources.
// <p>
// PlayerEndpoint will access the given resource, read all available data, and
// inject it into Kurento. Once this is is done, the injected video or audio
// will be available for passing through any other Filter or Endpoint to which
// the PlayerEndpoint gets connected.
// </p>
// <p>
// The source can provide either seekable or non-seekable media; this will
// dictate whether the PlayerEndpoint is able (or not) to seek through the file,
// for example to jump to any given timestamp.
// </p>
// <p>The <strong>Source URI</strong> supports these formats:</p>
// <ul>
// <li>
// File: A file path that will be read from the local file system. Example:
// <ul>
// <li><code>file:///path/to/file</code></li>
// </ul>
// </li>
// <li>
// HTTP: Any file available in an HTTP server. Examples:
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
// RTSP: Typically used to capture a feed from an IP Camera. Examples:
// <ul>
// <li><code>rtsp://{server-ip}</code></li>
// <li>
// <code>
// rtsp://{username}:{password}@{server-ip}:{server-port}/path/to/file
// </code>
// </li>
// </ul>
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
// <code>{username}</code> or <code>{password}</code> must be url-encoded.
// </strong>
// This means that colons (<code>:</code>) should be replaced with
// <code>%3A</code>, and 'at' signs (<code>@</code>) should be replaced with
// <code>%40</code>.
// </li>
// </ul>
// <p>
// Note that
// <strong> PlayerEndpoint requires read permissions to the source </strong>
// ; otherwise, the media server won't be able to retrieve any data, and an
// :rom:evt:`Error` will be fired. Make sure your application subscribes to this
// event, otherwise troubleshooting issues will be difficult.
// </p>
//
// <p>The list of valid operations is:</p>
// <ul>
// <li>
// <strong><code>play</code></strong>
// : Starts streaming media. If invoked after pause, it will resume playback.
// </li>
// <li>
// <strong><code>stop</code></strong>
// : Stops streaming media. If play is invoked afterwards, the file will be
// streamed from the beginning.
// </li>
// <li>
// <strong><code>pause</code></strong>
// : Pauses media streaming. Play must be invoked in order to resume playback.
// </li>
// <li>
// <strong><code>seek</code></strong>
// : If the source supports seeking to a different time position, then the
// PlayerEndpoint can:
// <ul>
// <li>
// <strong><code>setPosition</code></strong>
// : Allows to set the position in the file.
// </li>
// <li>
// <strong><code>getPosition</code></strong>
// : Returns the current position being streamed.
// </li>
// </ul>
// </li>
// </ul>
// <h2>Events fired</h2>
// <ul>
// <li>
// <strong>EndOfStreamEvent</strong>: If the file is streamed completely.
// </li>
// </ul>
//
type PlayerEndpoint struct {
	UriEndpoint

	// Returns info about the source being played
	VideoInfo *VideoInfo

	// Returns the GStreamer DOT string for this element's private pipeline
	ElementGstreamerDot string

	// Get or set the actual position of the video in ms. .. note:: Setting the position only works for seekable videos
	Position int64
}

// Return contructor params to be called by "Create".
func (elem *PlayerEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":   fmt.Sprintf("%s", from),
		"uri":             "",
		"useEncodedMedia": false,
		"networkCache":    2000,
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Starts reproducing the media, sending it to the `MediaSource`. If the endpoint
//
// has been connected to other endpoints, those will start receiving media.
func (elem *PlayerEndpoint) Play() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "play",
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
