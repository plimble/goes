package goes

// Encoder encode and decode event data
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, vPtr interface{}) error
}
