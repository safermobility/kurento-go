package kurento

type ElementConnectionData struct {
	Source            MediaElement
	Sink              MediaElement
	Type              MediaType
	SourceDescription string
	SinkDescription   string
}
