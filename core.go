package kurento

import "fmt"

// Base interface used to manage capabilities common to all Kurento elements.
// <h4>Properties</h4>
// <ul>
// <li>
// <b>id</b>: unique identifier assigned to this <code>MediaObject</code> at
// instantiation time. `MediaPipeline` IDs are generated with a GUID
// followed by suffix <code>_kurento.MediaPipeline</code>.
// `MediaElement` IDs are also a GUID with suffix
// <code>_kurento.{ElementType}</code> and prefixed by parent's ID.
// <blockquote>
// <dl>
// <dt><i>MediaPipeline ID example</i></dt>
// <dd>
// <code>
// 907cac3a-809a-4bbe-a93e-ae7e944c5cae_kurento.MediaPipeline
// </code>
// </dd>
// <dt><i>MediaElement ID example</i></dt>
// <dd>
// <code>
// 907cac3a-809a-4bbe-a93e-ae7e944c5cae_kurento.MediaPipeline/403da25a-805b-4cf1-8c55-f190588e6c9b_kurento.WebRtcEndpoint
// </code>
// </dd>
// </dl>
// </blockquote>
// </li>
// <li>
// <b>name</b>: free text intended to provide a friendly name for this
// <code>MediaObject</code>. Its default value is the same as the ID.
// </li>
// <li>
// <b>tags</b>: key-value pairs intended for applications to associate metadata
// to this <code>MediaObject</code> instance.
// </li>
// </ul>
// <p></p>
// <h4>Events</h4>
// <ul>
// <li>
// <strong>:rom:evt:`Error`<strong>: reports asynchronous error events. It is recommended to
// always subscribe a listener to this event, as regular error from the
// pipeline will be notified through it, instead of through an exception when
// invoking a method.
// </li>
// </ul>
//
type MediaObject struct {
	connection *Connection

	// `MediaPipeline` to which this <code>MediaObject</code> belongs. It returns itself when invoked for a pipeline object.
	MediaPipeline IMediaPipeline

	// Parent of this <code>MediaObject</code>.
	// <p>
	// The parent of a `Hub` or a `MediaElement` is its
	// `MediaPipeline`. A `MediaPipeline` has no parent, so this
	// property will be null.
	// </p>
	//
	Parent IMediaObject

	// Unique identifier of this <code>MediaObject</code>.
	// <p>
	// It's a synthetic identifier composed by a GUID and
	// <code>MediaObject</code> type. The ID is prefixed with the parent ID when the
	// object has parent: <i>ID_parent/ID_media-object</i>.
	// </p>
	//
	Id string

	// Children of this <code>MediaObject</code>.
	// @deprecated Use children instead.
	//
	Childs []IMediaObject

	// Children of this <code>MediaObject</code>.
	Children []IMediaObject

	// This <code>MediaObject</code>'s name.
	// <p>
	// This is just sugar to simplify developers' life debugging, it is not used
	// internally for indexing nor identifying the objects. By default, it's the
	// object's ID.
	// </p>
	//
	Name string

	// Flag activating or deactivating sending the element's tags in fired events.
	SendTagsInEvents bool

	// <code>MediaObject</code> creation time in seconds since Epoch.
	CreationTime int
}

// Return contructor params to be called by "Create".
func (elem *MediaObject) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Adds a new tag to this <code>MediaObject</code>.
// If the tag is already present, it changes the value.
//
func (elem *MediaObject) AddTag(key string, value string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)
	setIfNotEmpty(params, "value", value)

	reqparams := map[string]interface{}{
		"operation":       "addTag",
		"object":          elem.Id,
		"operationParams": params,
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

// Removes an existing tag.
// Exists silently with no error if tag is not defined.
//
func (elem *MediaObject) RemoveTag(key string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)

	reqparams := map[string]interface{}{
		"operation":       "removeTag",
		"object":          elem.Id,
		"operationParams": params,
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

// Returns the value of given tag, or MEDIA_OBJECT_TAG_KEY_NOT_FOUND if tag is not defined.
// Returns:
// // The value associated to the given key.
func (elem *MediaObject) GetTag(key string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)

	reqparams := map[string]interface{}{
		"operation":       "getTag",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The value associated to the given key.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

// Returns all tags attached to this <code>MediaObject</code>.
// Returns:
// // An array containing all key-value pairs associated with this <code>MediaObject</code>.
func (elem *MediaObject) GetTags() ([]Tag, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getTags",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // An array containing all key-value pairs associated with this <code>MediaObject</code>.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	ret := []Tag{}
	return ret, err

}

type IServerManager interface {
	GetKmd(moduleName string) (string, error)
	GetCpuCount() (int, error)
	GetUsedCpu(interval int) (float64, error)
	GetUsedMemory() (int64, error)
}

// This is a standalone object for managing the MediaServer
type ServerManager struct {
	MediaObject

	// Server information, version, modules, factories, etc
	Info *ServerInfo

	// All the pipelines available in the server
	Pipelines []IMediaPipeline

	// All active sessions in the server
	Sessions []string

	// Metadata stored in the server
	Metadata string
}

// Return contructor params to be called by "Create".
func (elem *ServerManager) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Returns the kmd associated to a module
// Returns:
// // The kmd file.
func (elem *ServerManager) GetKmd(moduleName string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "moduleName", moduleName)

	reqparams := map[string]interface{}{
		"operation":       "getKmd",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The kmd file.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

// Number of CPU cores that the media server can use.
// <p>
// Linux processes can be configured to use only a subset of the cores that are
// available in the system, via the process affinity settings
// (<strong>sched_setaffinity(2)</strong>). With this method it is possible to
// know the number of cores that the media server can use in the machine where it
// is running.
// </p>
// <p>
// For example, it's possible to limit the core affinity inside a Docker
// container by running with a command such as
// <em>docker run --cpuset-cpus='0,1'</em>.
// </p>
// <p>
// Note that the return value represents the number of
// <em>logical</em> processing units available, i.e. CPU cores including
// Hyper-Threading.
// </p>
//
// Returns:
// // Number of CPU cores available for the media server.
func (elem *ServerManager) GetCpuCount() (int, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getCpuCount",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // Number of CPU cores available for the media server.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(int); ok {
		return value, err
	}

	return 0, err

}

// Average CPU usage of the server.
// <p>
// This method measures the average CPU usage of the media server during the
// requested interval. Normally you will want to choose an interval between 1000
// and 10000 ms.
// </p>
// <p>
// The returned value represents the global system CPU usage of the media server,
// as an average across all processing units (CPU cores).
// </p>
//
// Returns:
// // CPU usage %.
func (elem *ServerManager) GetUsedCpu(interval int) (float64, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "interval", interval)

	reqparams := map[string]interface{}{
		"operation":       "getUsedCpu",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // CPU usage %.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(float64); ok {
		return value, err
	}

	return 0, err

}

// Returns the amount of memory that the server is using, in KiB
// Returns:
// // Used memory, in KiB.
func (elem *ServerManager) GetUsedMemory() (int64, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getUsedMemory",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // Used memory, in KiB.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(int64); ok {
		return value, err
	}

	return 0, err

}

type ISessionEndpoint interface {
}

// All networked Endpoints that require to manage connection sessions with remote peers implement this interface.
type SessionEndpoint struct {
	Endpoint
}

// Return contructor params to be called by "Create".
func (elem *SessionEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHub interface {
	GetGstreamerDot(details GstreamerDotDetails) (string, error)
}

// A Hub is a routing `MediaObject`.
// It connects several `endpoints <Endpoint>` together
//
type Hub struct {
	MediaObject
}

// Return contructor params to be called by "Create".
func (elem *Hub) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Returns a string in dot (graphviz) format that represents the gstreamer elements inside the pipeline
// Returns:
// // The dot graph.
func (elem *Hub) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	reqparams := map[string]interface{}{
		"operation":       "getGstreamerDot",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The dot graph.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

type IFilter interface {
}

// Base interface for all filters.
// <p>
// This is a certain type of `MediaElement`, that processes media
// injected through its sinks, and delivers the outcome through its sources.
// </p>
//
type Filter struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *Filter) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IEndpoint interface {
}

// Base interface for all end points.
// <p>
// An Endpoint is a `MediaElement` that allows Kurento to exchange
// media contents with external systems, supporting different transport protocols
// and mechanisms, such as RTP, WebRTC, HTTP(s), "file://" URLs, etc.
// </p>
// <p>
// An "Endpoint" may contain both sources and sinks for different media types,
// to provide bidirectional communication.
// </p>
//
type Endpoint struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *Endpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHubPort interface {
}

// This `MediaElement` specifies a connection with a `Hub`
type HubPort struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *HubPort) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"hub": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IPassThrough interface {
}

// This `MediaElement` that just passes media through
type PassThrough struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *PassThrough) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IUriEndpoint interface {
	Pause() error
	Stop() error
}

// Interface for endpoints the require a URI to work.
// An example of this, would be a `PlayerEndpoint` whose URI property could be used to locate a file to stream.
//
type UriEndpoint struct {
	Endpoint

	// The uri for this endpoint.
	Uri string

	// State of the endpoint
	State *UriEndpointState
}

// Return contructor params to be called by "Create".
func (elem *UriEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Pauses the feed
func (elem *UriEndpoint) Pause() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "pause",
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

// Stops the feed
func (elem *UriEndpoint) Stop() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "stop",
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

type IMediaPipeline interface {
	GetGstreamerDot(details GstreamerDotDetails) (string, error)
}

// A pipeline is a container for a collection of `MediaElements<MediaElement>` and `MediaMixers<MediaMixer>`.
// It offers the methods needed to control the creation and connection of elements inside a certain pipeline.
//
type MediaPipeline struct {
	MediaObject

	// If statistics about pipeline latency are enabled for all mediaElements
	LatencyStats bool
}

// Return contructor params to be called by "Create".
func (elem *MediaPipeline) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Returns a string in dot (graphviz) format that represents the gstreamer elements inside the pipeline
// Returns:
// // The dot graph.
func (elem *MediaPipeline) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	reqparams := map[string]interface{}{
		"operation":       "getGstreamerDot",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The dot graph.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

type ISdpEndpoint interface {
	GenerateOffer(options OfferOptions) (string, error)
	ProcessOffer(offer string) (string, error)
	ProcessAnswer(answer string) (string, error)
	GetLocalSessionDescriptor() (string, error)
	GetRemoteSessionDescriptor() (string, error)
}

// Interface implemented by Endpoints that require an SDP Offer/Answer negotiation in order to configure a media session.
// <p>Functionality provided by this API:</p>
// <ul>
// <li>Generate SDP offers.</li>
// <li>Process SDP offers.</li>
// <li>Configure SDP related params.</li>
// </ul>
//
type SdpEndpoint struct {
	SessionEndpoint

	// Maximum input bitrate, signaled in SDP Offers to WebRTC and RTP senders.
	// <p>
	// This is used to put a limit on the bitrate that the remote peer will send to
	// this endpoint. The net effect of setting this parameter is that
	// <i>when Kurento generates an SDP Offer</i>, an 'Application Specific' (AS)
	// maximum bandwidth attribute will be added to the SDP media section:
	// <code>b=AS:{value}</code>.
	// </p>
	// <p>Note: This parameter has to be set before the SDP is generated.</p>
	// <ul>
	// <li>Unit: kbps (kilobits per second).</li>
	// <li>Default: 0.</li>
	// <li>0 = unlimited.</li>
	// </ul>
	//
	MaxAudioRecvBandwidth int

	// Maximum input bitrate, signaled in SDP Offers to WebRTC and RTP senders.
	// <p>
	// This is used to put a limit on the bitrate that the remote peer will send to
	// this endpoint. The net effect of setting this parameter is that
	// <i>when Kurento generates an SDP Offer</i>, an 'Application Specific' (AS)
	// maximum bandwidth attribute will be added to the SDP media section:
	// <code>b=AS:{value}</code>.
	// </p>
	// <p>Note: This parameter has to be set before the SDP is generated.</p>
	// <ul>
	// <li>Unit: kbps (kilobits per second).</li>
	// <li>Default: 0.</li>
	// <li>0 = unlimited.</li>
	// </ul>
	//
	MaxVideoRecvBandwidth int
}

// Return contructor params to be called by "Create".
func (elem *SdpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Generates an SDP offer with media capabilities of the Endpoint.
// Throws:
// <ul>
// <li>
// SDP_END_POINT_ALREADY_NEGOTIATED If the endpoint is already negotiated.
// </li>
// <li>
// SDP_END_POINT_GENERATE_OFFER_ERROR if the generated offer is empty. This is
// most likely due to an internal error.
// </li>
// </ul>
//
// Returns:
// // The SDP offer.
func (elem *SdpEndpoint) GenerateOffer(options OfferOptions) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "options", options)

	reqparams := map[string]interface{}{
		"operation":       "generateOffer",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The SDP offer.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

// Processes SDP offer of the remote peer, and generates an SDP answer based on the endpoint's capabilities.
// <p>
// If no matching capabilities are found, the SDP will contain no codecs.
// </p>
// Throws:
// <ul>
// <li>
// SDP_PARSE_ERROR If the offer is empty or has errors.
// </li>
// <li>
// SDP_END_POINT_ALREADY_NEGOTIATED If the endpoint is already negotiated.
// </li>
// <li>
// SDP_END_POINT_PROCESS_OFFER_ERROR if the generated offer is empty. This is
// most likely due to an internal error.
// </li>
// </ul>
//
// Returns:
// // The chosen configuration from the ones stated in the SDP offer.
func (elem *SdpEndpoint) ProcessOffer(offer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "offer", offer)

	reqparams := map[string]interface{}{
		"operation":       "processOffer",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The chosen configuration from the ones stated in the SDP offer.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

// Generates an SDP offer with media capabilities of the Endpoint.
// Throws:
// <ul>
// <li>
// SDP_PARSE_ERROR If the offer is empty or has errors.
// </li>
// <li>
// SDP_END_POINT_ALREADY_NEGOTIATED If the endpoint is already negotiated.
// </li>
// <li>
// SDP_END_POINT_PROCESS_ANSWER_ERROR if the result of processing the answer is
// an empty string. This is most likely due to an internal error.
// </li>
// <li>
// SDP_END_POINT_NOT_OFFER_GENERATED If the method is invoked before the
// generateOffer method.
// </li>
// </ul>
//
// Returns:
// // Updated SDP offer, based on the answer received.
func (elem *SdpEndpoint) ProcessAnswer(answer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "answer", answer)

	reqparams := map[string]interface{}{
		"operation":       "processAnswer",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // Updated SDP offer, based on the answer received.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

// Returns the local SDP.
// <ul>
// <li>
// No offer has been generated: returns null.
// </li>
// <li>
// Offer has been generated: returns the SDP offer.
// </li>
// <li>
// Offer has been generated and answer processed: returns the agreed SDP.
// </li>
// </ul>
//
// Returns:
// // The last agreed SessionSpec.
func (elem *SdpEndpoint) GetLocalSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getLocalSessionDescriptor",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The last agreed SessionSpec.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

// This method returns the remote SDP.
// If the negotiation process is not complete, it will return NULL.
//
// Returns:
// // The last agreed User Agent session description.
func (elem *SdpEndpoint) GetRemoteSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getRemoteSessionDescriptor",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The last agreed User Agent session description.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

type IBaseRtpEndpoint interface {
}

// Handles RTP communications.
// <p>
// All endpoints that rely on the RTP protocol, like the
// <strong>RtpEndpoint</strong> or the <strong>WebRtcEndpoint</strong>, inherit
// from this class. The endpoint provides information about the Connection state
// and the Media state, which can be consulted at any time through the
// :rom:attr:`getMediaState` and the :rom:attr:`getConnectionState` methods.
// It is also possible subscribe to events fired when these properties change:
// </p>
// <ul>
// <li>
// <strong>:rom:evt:`ConnectionStateChanged`</strong>: This event is raised
// when the connection between two peers changes. It can have two values:
// <ul>
// <li>CONNECTED</li>
// <li>DISCONNECTED</li>
// </ul>
// </li>
// <li>
// <strong>:rom:evt:`MediaStateChanged`</strong>: This event provides
// information about the state of the underlying RTP session. Possible values
// are:
// <ul>
// <li>CONNECTED: There is an RTCP packet flow between peers.</li>
// <li>
// DISCONNECTED: Either no RTCP packets have been received yet, or the
// remote peer has ended the RTP session with a <code>BYE</code> message,
// or at least 5 seconds have elapsed since the last RTCP packet was
// received.
// </li>
// </ul>
// <p>
// The standard definition of RTP (<a
// href='https://tools.ietf.org/html/rfc3550'
// target='_blank'
// >RFC 3550</a
// >) describes a session as active whenever there is a maintained flow of
// RTCP control packets, regardless of whether there is actual media flowing
// through RTP data packets or not. The reasoning behind this is that, at any
// given moment, a participant of an RTP session might temporarily stop
// sending RTP data packets, but this wouldn't necessarily mean that the RTP
// session as a whole is finished; it maybe just means that the participant
// has some temporary issues but it will soon resume sending data. For this
// reason, that an RTP session has really finished is something that is
// considered only by the prolonged absence of RTCP control packets between
// participants.
// </p>
// <p>
// Since RTCP packets do not flow at a constant rate (for instance,
// minimizing a browser window with a WebRTC's
// <code>RTCPeerConnection</code> object might affect the sending interval),
// it is not possible to immediately detect their absence and assume that the
// RTP session has finished. Instead, there is a guard period of
// approximately <strong>5 seconds</strong> of missing RTCP packets before
// considering that the underlying RTP session is effectively finished, thus
// triggering a <code>MediaStateChangedEvent = DISCONNECTED</code> event.
// </p>
// <p>
// In other words, there is always a period during which there might be no
// media flowing, but this event hasn't been fired yet. Nevertheless, this is
// the most reliable and useful way of knowing what is the long-term, steady
// state of RTP media exchange.
// </p>
// <p>
// The :rom:evt:`ConnectionStateChanged` comes in contrast with more
// instantaneous events such as MediaElement's
// :rom:evt:`MediaFlowInStateChanged` and
// :rom:evt:`MediaFlowOutStateChanged`, which are triggered almost
// immediately after the RTP data packets stop flowing between RTP session
// participants. This makes the <em>MediaFlow</em> events a good way to
// know if participants are suffering from short-term intermittent
// connectivity issues, but they are not enough to know if the connectivity
// issues are just spurious network hiccups or are part of a more long-term
// disconnection problem.
// </p>
// </li>
// </ul>
// <p>
// Part of the bandwidth control for the video component of the media session is
// done here:
// </p>
// <ul>
// <li>
// Input bandwidth: Values used to inform remote peers about the bitrate that
// can be sent to this endpoint.
// <ul>
// <li>
// <strong>MinVideoRecvBandwidth</strong>: Minimum input bitrate, requested
// from WebRTC senders with REMB (Default: 30 Kbps).
// </li>
// <li>
// <strong>MaxAudioRecvBandwidth</strong> and
// <strong>MaxVideoRecvBandwidth</strong>: Maximum input bitrate, signaled
// in SDP Offers to WebRTC and RTP senders (Default: unlimited).
// </li>
// </ul>
// </li>
// <li>
// Output bandwidth: Values used to control bitrate of the video streams sent
// to remote peers. It is important to keep in mind that pushed bitrate depends
// on network and remote peer capabilities. Remote peers can also announce
// bandwidth limitation in their SDPs (through the
// <code>b={modifier}:{value}</code> attribute). Kurento will always enforce
// bitrate limitations specified by the remote peer over internal
// configurations.
// <ul>
// <li>
// <strong>MinVideoSendBandwidth</strong>: REMB override of minimum bitrate
// sent to WebRTC receivers (Default: 100 Kbps).
// </li>
// <li>
// <strong>MaxVideoSendBandwidth</strong>: REMB override of maximum bitrate
// sent to WebRTC receivers (Default: 500 Kbps).
// </li>
// <li>
// <strong>RembParams.rembOnConnect</strong>: Initial local REMB bandwidth
// estimation that gets propagated when a new endpoint is connected.
// </li>
// </ul>
// </li>
// </ul>
// <p>
// <strong>
// All bandwidth control parameters must be changed before the SDP negotiation
// takes place, and can't be changed afterwards.
// </strong>
// </p>
//
type BaseRtpEndpoint struct {
	SdpEndpoint

	// Minimum input bitrate, requested from WebRTC senders with REMB.
	// <p>
	// This is used to set a minimum value of local REMB during bandwidth estimation,
	// if supported by the implementing class. The REMB estimation will then be sent
	// to remote peers, requesting them to send at least the indicated video bitrate.
	// It follows that min values will only have effect in remote peers that support
	// this congestion control mechanism, such as Chrome.
	// </p>
	// <ul>
	// <li>Unit: kbps (kilobits per second).</li>
	// <li>Default: 0.</li>
	// <li>
	// Note: The absolute minimum REMB value is 30 kbps, even if a lower value is
	// set here.
	// </li>
	// </ul>
	//
	MinVideoRecvBandwidth int

	// REMB override of minimum bitrate sent to WebRTC receivers.
	// <p>
	// With this parameter you can control the minimum video quality that will be
	// sent when reacting to bad network conditions. Setting this parameter to a low
	// value permits the video quality to drop when the network conditions get worse.
	// </p>
	// <p>
	// This parameter provides a way to override the bitrate requested by remote REMB
	// bandwidth estimations: the bitrate sent will be always equal or greater than
	// this parameter, even if the remote peer requests even lower bitrates.
	// </p>
	// <p>
	// Note that if you set this parameter too high (trying to avoid bad video
	// quality altogether), you would be limiting the adaptation ability of the
	// congestion control algorithm, and your stream might be unable to ever recover
	// from adverse network conditions.
	// </p>
	// <ul>
	// <li>Unit: kbps (kilobits per second).</li>
	// <li>Default: 100.</li>
	// <li>
	// 0 = unlimited: the video bitrate will drop as needed, even to the lowest
	// possible quality, which might make the video completely blurry and
	// pixelated.
	// </li>
	// </ul>
	//
	MinVideoSendBandwidth int

	// REMB override of maximum bitrate sent to WebRTC receivers.
	// <p>
	// With this parameter you can control the maximum video quality that will be
	// sent when reacting to good network conditions. Setting this parameter to a
	// high value permits the video quality to raise when the network conditions get
	// better.
	// </p>
	// <p>
	// This parameter provides a way to limit the bitrate requested by remote REMB
	// bandwidth estimations: the bitrate sent will be always equal or less than this
	// parameter, even if the remote peer requests higher bitrates.
	// </p>
	// <p>
	// Note that the default value of <strong>500 kbps</strong> is a VERY
	// conservative one, and leads to a low maximum video quality. Most applications
	// will probably want to increase this to higher values such as 2000 kbps (2
	// mbps).
	// </p>
	// <p>
	// The REMB congestion control algorithm works by gradually increasing the output
	// video bitrate, until the available bandwidth is fully used or the maximum send
	// bitrate has been reached. This is a slow, progressive change, which starts at
	// 300 kbps by default. You can change the default starting point of REMB
	// estimations, by setting <code>RembParams.rembOnConnect</code>.
	// </p>
	// <ul>
	// <li>Unit: kbps (kilobits per second).</li>
	// <li>Default: 500.</li>
	// <li>
	// 0 = unlimited: the video bitrate will grow until all the available network
	// bandwidth is used by the stream.<br />
	// Note that this might have a bad effect if more than one stream is running
	// (as all of them would try to raise the video bitrate indefinitely, until the
	// network gets saturated).
	// </li>
	// </ul>
	//
	MaxVideoSendBandwidth int

	// Media flow state.
	// <ul>
	// <li>CONNECTED: There is an RTCP flow.</li>
	// <li>DISCONNECTED: No RTCP packets have been received for at least 5 sec.</li>
	// </ul>
	//
	MediaState *MediaState

	// Connection state.
	// <ul>
	// <li>CONNECTED</li>
	// <li>DISCONNECTED</li>
	// </ul>
	//
	ConnectionState *ConnectionState

	// Maximum Transmission Unit (MTU) used for RTP.
	// <p>
	// This setting affects the maximum size that will be used by RTP payloads. You
	// can change it from the default, if you think that a different value would be
	// beneficial for the typical network settings of your application.
	// </p>
	// <p>
	// The default value is 1200 Bytes. This is the same as in <b>libwebrtc</b> (from
	// webrtc.org), as used by
	// <a
	// href='https://dxr.mozilla.org/mozilla-central/rev/b5c5ba07d3dbd0d07b66fa42a103f4df2c27d3a2/media/webrtc/trunk/webrtc/media/engine/constants.cc#16'
	// >Firefox</a
	// >
	// or
	// <a
	// href='https://source.chromium.org/chromium/external/webrtc/src/+/6dd488b2e55125644263e4837f1abd950d5e410d:media/engine/constants.cc;l=15'
	// >Chrome</a
	// >
	// . You can read more about this value in
	// <a
	// href='https://groups.google.com/d/topic/discuss-webrtc/gH5ysR3SoZI/discussion'
	// >Why RTP max packet size is 1200 in WebRTC?</a
	// >
	// .
	// </p>
	// <p>
	// <b>WARNING</b>: Change this value ONLY if you really know what you are doing
	// and you have strong reasons to do so. Do NOT change this parameter just
	// because it <i>seems</i> to work better for some reduced scope tests. The
	// default value is a consensus chosen by people who have deep knowledge about
	// network optimization.
	// </p>
	// <ul>
	// <li>Unit: Bytes.</li>
	// <li>Default: 1200.</li>
	// </ul>
	//
	Mtu int

	// Advanced parameters to configure the congestion control algorithm.
	RembParams *RembParams
}

// Return contructor params to be called by "Create".
func (elem *BaseRtpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IMediaElement interface {
	GetSourceConnections(mediaType MediaType, description string) ([]ElementConnectionData, error)
	GetSinkConnections(mediaType MediaType, description string) ([]ElementConnectionData, error)
	Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error
	Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error
	SetAudioFormat(caps AudioCaps) error
	SetVideoFormat(caps VideoCaps) error
	GetGstreamerDot(details GstreamerDotDetails) (string, error)
	SetOutputBitrate(bitrate int) error
	GetStats(mediaType MediaType) (map[string]Stats, error)
	IsMediaFlowingIn(mediaType MediaType, sinkMediaDescription string) (bool, error)
	IsMediaFlowingOut(mediaType MediaType, sourceMediaDescription string) (bool, error)
	IsMediaTranscoding(mediaType MediaType, binName string) (bool, error)
}

// The basic building block of the media server, that can be interconnected inside a pipeline.
// <p>
// A `MediaElement` is a module that encapsulates a specific media
// capability, and that is able to exchange media with other MediaElements
// through an internal element called <b>pad</b>.
// </p>
// <p>
// A pad can be defined as an input or output interface. Input pads are called
// sinks, and it's where the media elements receive media from other media
// elements. Output interfaces are called sources, and it's the pad used by the
// media element to feed media to other media elements. There can be only one
// sink pad per media element. On the other hand, the number of source pads is
// unconstrained. This means that a certain media element can receive media only
// from one element at a time, while it can send media to many others. Pads are
// created on demand, when the connect method is invoked. When two media elements
// are connected, one media pad is created for each type of media connected. For
// example, if you connect AUDIO and VIDEO between two media elements, each one
// will need to create two new pads: one for AUDIO and one for VIDEO.
// </p>
// <p>
// When media elements are connected, it can be the case that the encoding
// required in both input and output pads is not the same, and thus it needs to
// be transcoded. This is something that is handled transparently by the
// MediaElement internals, but such transcoding has a toll in the form of a
// higher CPU load, so connecting MediaElements that need media encoded in
// different formats is something to consider as a high load operation. The event
// `MediaTranscodingStateChanged` allows to inform the client application of
// whether media transcoding is being enabled or not inside any MediaElement
// object.
// </p>
//
type MediaElement struct {
	MediaObject

	// Minimum video bandwidth for transcoding.
	// @deprecated Deprecated due to a typo. Use :rom:meth:`minOutputBitrate` instead of this function.
	MinOuputBitrate int

	// Minimum video bitrate for transcoding.
	// <ul>
	// <li>Unit: bps (bits per second).</li>
	// <li>Default: 0.</li>
	// </ul>
	//
	MinOutputBitrate int

	// Maximum video bandwidth for transcoding.
	// @deprecated Deprecated due to a typo. Use :rom:meth:`maxOutputBitrate` instead of this function.
	MaxOuputBitrate int

	// Maximum video bitrate for transcoding.
	// <ul>
	// <li>Unit: bps (bits per second).</li>
	// <li>Default: MAXINT.</li>
	// <li>0 = unlimited.</li>
	// </ul>
	//
	MaxOutputBitrate int
}

// Return contructor params to be called by "Create".
func (elem *MediaElement) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Gets information about the sink pads of this media element.
// <p>
// Since sink pads are the interface through which a media element gets it's
// media, whatever is connected to an element's sink pad is formally a source of
// media. Media can be filtered by type, or by the description given to the pad
// though which both elements are connected.
// </p>
//
// Returns:
// // A list of the connections information that are sending media to this element. The list will be empty if no sources are found.
func (elem *MediaElement) GetSourceConnections(mediaType MediaType, description string) ([]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "description", description)

	reqparams := map[string]interface{}{
		"operation":       "getSourceConnections",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // A list of the connections information that are sending media to this element. The list will be empty if no sources are found.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	ret := []ElementConnectionData{}
	return ret, err

}

// Gets information about the source pads of this media element.
// <p>
// Since source pads connect to other media element's sinks, this is formally the
// sink of media from the element's perspective. Media can be filtered by type,
// or by the description given to the pad though which both elements are
// connected.
// </p>
//
// Returns:
// // A list of the connections information that are receiving media from this element. The list will be empty if no sources are found.
func (elem *MediaElement) GetSinkConnections(mediaType MediaType, description string) ([]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "description", description)

	reqparams := map[string]interface{}{
		"operation":       "getSinkConnections",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // A list of the connections information that are receiving media from this element. The list will be empty if no sources are found.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	ret := []ElementConnectionData{}
	return ret, err

}

// Connects two elements, with the media flowing from left to right.
// <p>
// The element that invokes the connect will be the source of media, creating one
// sink pad for each type of media connected. The element given as parameter to
// the method will be the sink, and it will create one sink pad per media type
// connected.
// </p>
// <p>
// If otherwise not specified, all types of media are connected by default
// (AUDIO, VIDEO and DATA). It is recommended to connect the specific types of
// media if not all of them will be used. For this purpose, the connect method
// can be invoked more than once on the same two elements, but with different
// media types.
// </p>
// <p>
// The connection is unidirectional. If a bidirectional connection is desired,
// the position of the media elements must be inverted. For instance,
// webrtc1.connect(webrtc2) is connecting webrtc1 as source of webrtc2. In order
// to create a WebRTC one-2one conversation, the user would need to specify the
// connection on the other direction with webrtc2.connect(webrtc1).
// </p>
// <p>
// Even though one media element can have one sink pad per type of media, only
// one media element can be connected to another at a given time. If a media
// element is connected to another, the former will become the source of the sink
// media element, regardless whether there was another element connected or not.
// </p>
//
func (elem *MediaElement) Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)
	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)
	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	reqparams := map[string]interface{}{
		"operation":       "connect",
		"object":          elem.Id,
		"operationParams": params,
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

// Disconnects two media elements. This will release the source pads of the source media element, and the sink pads of the sink media element.
func (elem *MediaElement) Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)
	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)
	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	reqparams := map[string]interface{}{
		"operation":       "disconnect",
		"object":          elem.Id,
		"operationParams": params,
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

// Set the type of data for the audio stream.
// <p>
// MediaElements that do not support configuration of audio capabilities will
// throw a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR exception.
// </p>
// <p>
// NOTE: This method is not implemented yet by the Media Server to do anything
// useful.
// </p>
//
func (elem *MediaElement) SetAudioFormat(caps AudioCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	reqparams := map[string]interface{}{
		"operation":       "setAudioFormat",
		"object":          elem.Id,
		"operationParams": params,
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

// Set the type of data for the video stream.
// <p>
// MediaElements that do not support configuration of video capabilities will
// throw a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR exception
// </p>
// <p>
// NOTE: This method is not implemented yet by the Media Server to do anything
// useful.
// </p>
//
func (elem *MediaElement) SetVideoFormat(caps VideoCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	reqparams := map[string]interface{}{
		"operation":       "setVideoFormat",
		"object":          elem.Id,
		"operationParams": params,
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

// Return a .dot file describing the topology of the media element.
// <p>The element can be queried for certain type of data:</p>
// <ul>
// <li>SHOW_ALL: default value</li>
// <li>SHOW_CAPS_DETAILS</li>
// <li>SHOW_FULL_PARAMS</li>
// <li>SHOW_MEDIA_TYPE</li>
// <li>SHOW_NON_DEFAULT_PARAMS</li>
// <li>SHOW_STATES</li>
// <li>SHOW_VERBOSE</li>
// </ul>
//
// Returns:
// // The dot graph.
func (elem *MediaElement) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	reqparams := map[string]interface{}{
		"operation":       "getGstreamerDot",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The dot graph.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(string); ok {
		return value, err
	}

	return "", err

}

// @deprecated
// Allows change the target bitrate for the media output, if the media is encoded using VP8 or H264. This method only works if it is called before the media starts to flow.
func (elem *MediaElement) SetOutputBitrate(bitrate int) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "bitrate", bitrate)

	reqparams := map[string]interface{}{
		"operation":       "setOutputBitrate",
		"object":          elem.Id,
		"operationParams": params,
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

// Gets the statistics related to an endpoint. If no media type is specified, it returns statistics for all available types.
// Returns:
// // Delivers a successful result in the form of a RTC stats report. A RTC stats report represents a map between strings, identifying the inspected objects (RTCStats.id), and their corresponding RTCStats objects.
func (elem *MediaElement) GetStats(mediaType MediaType) (map[string]Stats, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)

	reqparams := map[string]interface{}{
		"operation":       "getStats",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // Delivers a successful result in the form of a RTC stats report. A RTC stats report represents a map between strings, identifying the inspected objects (RTCStats.id), and their corresponding RTCStats objects.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	ret := map[string]Stats{}
	return ret, err

}

// This method indicates whether the media element is receiving media of a certain type. The media sink pad can be identified individually, if needed. It is only supported for AUDIO and VIDEO types, raising a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR otherwise. If the pad indicated does not exist, if will return false.
// Returns:
// // TRUE if there is media, FALSE in other case.
func (elem *MediaElement) IsMediaFlowingIn(mediaType MediaType, sinkMediaDescription string) (bool, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	reqparams := map[string]interface{}{
		"operation":       "isMediaFlowingIn",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // TRUE if there is media, FALSE in other case.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(bool); ok {
		return value, err
	}

	return false, err

}

// This method indicates whether the media element is emitting media of a certain type. The media source pad can be identified individually, if needed. It is only supported for AUDIO and VIDEO types, raising a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR otherwise. If the pad indicated does not exist, if will return false.
// Returns:
// // TRUE if there is media, FALSE in other case.
func (elem *MediaElement) IsMediaFlowingOut(mediaType MediaType, sourceMediaDescription string) (bool, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)

	reqparams := map[string]interface{}{
		"operation":       "isMediaFlowingOut",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // TRUE if there is media, FALSE in other case.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(bool); ok {
		return value, err
	}

	return false, err

}

// Indicates whether this media element is actively transcoding between input and output pads. This operation is only supported for AUDIO and VIDEO media types, raising a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR otherwise.
// The internal GStreamer processing bin can be indicated, if needed; if the bin doesn't exist, the return value will be FALSE.
// Returns:
// // TRUE if media is being transcoded, FALSE otherwise.
func (elem *MediaElement) IsMediaTranscoding(mediaType MediaType, binName string) (bool, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "binName", binName)

	reqparams := map[string]interface{}{
		"operation":       "isMediaTranscoding",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // TRUE if media is being transcoded, FALSE otherwise.
	var err error
	if response.Error != nil {
		err = fmt.Errorf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	if value, ok := response.Result["value"].(bool); ok {
		return value, err
	}

	return false, err

}