package kurento

type ServerInfo struct {
	Version      string
	Modules      []ModuleInfo
	Type         ServerType
	Capabilities []string
}
