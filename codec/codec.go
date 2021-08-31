package codec

type MarshalerType string

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}
