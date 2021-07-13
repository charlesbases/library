package metadata

type Metadata map[string]interface{}

// NewMetadata .
func NewMetadata() Metadata {
	return make(map[string]interface{})
}
