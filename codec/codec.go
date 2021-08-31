package codec

type MarshalerType string

const (
	MarshalerType_Json  MarshalerType = "application/json"
	MarshalerType_Proto MarshalerType = "application/proto"
)

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() MarshalerType
}
