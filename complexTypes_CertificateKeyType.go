package kurento

// .
type CertificateKeyType string

// Implement fmt.Stringer interface
func (t CertificateKeyType) String() string {
	return string(t)
}

const (
	CERTIFICATEKEYTYPE_RSA   CertificateKeyType = "RSA"
	CERTIFICATEKEYTYPE_ECDSA CertificateKeyType = "ECDSA"
)