package jwt

// Alg is algorithm interface
//
type Alg interface {
	// Name return the alogrithm's name
	Name() string
	// Sign build data signature with alogrithm defined
	Sign(data []byte) ([]byte, error)
	// Verify check data signature
	Verify(data, sign []byte) error
}
