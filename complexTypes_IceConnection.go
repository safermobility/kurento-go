package kurento

type IceConnection struct {
	StreamId    string
	ComponentId int
	State       IceComponentState
}
