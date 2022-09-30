package kurento

// Indicates if the server is a real media server or a proxy
type ServerType string

// Implement fmt.Stringer interface
func (t ServerType) String() string {
	return string(t)
}

const (
	SERVERTYPE_KMS ServerType = "KMS"
	SERVERTYPE_KCS ServerType = "KCS"
)
