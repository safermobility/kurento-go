package kurento

// Describes the encryption and authentication algorithms
type CryptoSuite string

// Implement fmt.Stringer interface
func (t CryptoSuite) String() string {
	return string(t)
}

const (
	CRYPTOSUITE_AES_128_CM_HMAC_SHA1_32 CryptoSuite = "AES_128_CM_HMAC_SHA1_32"
	CRYPTOSUITE_AES_128_CM_HMAC_SHA1_80 CryptoSuite = "AES_128_CM_HMAC_SHA1_80"
	CRYPTOSUITE_AES_256_CM_HMAC_SHA1_32 CryptoSuite = "AES_256_CM_HMAC_SHA1_32"
	CRYPTOSUITE_AES_256_CM_HMAC_SHA1_80 CryptoSuite = "AES_256_CM_HMAC_SHA1_80"
)
