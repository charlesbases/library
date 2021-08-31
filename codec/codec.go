package codec

type MarshalerType string

const (
	MarshalerType_Json  = "json"
	MarshalerType_Proto = "proto"
)

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() MarshalerType
}
