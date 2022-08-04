package codec

type ContentType string

const (
	// ContentTypeJson json
	ContentTypeJson ContentType = "application/json"
	// ContentTypeProto proto
	ContentTypeProto ContentType = "application/proto"
)

// Marshaler 编解码器
type Marshaler interface {
	Marshal(interface{}, ...MarshalOption) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	ContentType() ContentType
}

// String .
func (ct ContentType) String() string {
	return string(ct)
}
