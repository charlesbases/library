package content

// Type content-type
type Type int8

// DefaultContentType default of content-type
const DefaultContentType Type = Json

const (
	// Text application/text
	Text Type = iota
	// Yaml application/yaml
	Yaml
	// Json application/json
	Json
	// Proto application/proto
	Proto
	// Bytes application/bytes
	Bytes
	// Stream application/octet-stream
	Stream
	// FromData multiparty/from-data
	FromData
	// Zip application/zip
	Zip
)

var contents = map[Type]string{
	Zip:      "application/zip",
	Yaml:     "application/yaml",
	Text:     "application/text",
	Json:     "application/json",
	Bytes:    "application/bytes",
	Proto:    "application/proto",
	Stream:   "application/octet-stream",
	FromData: "multiparty/from-data",
}

var reverse = map[string]Type{
	"application/zip":          Zip,
	"application/yaml":         Yaml,
	"application/text":         Text,
	"application/json":         Json,
	"application/bytes":        Bytes,
	"application/proto":        Proto,
	"application/octet-stream": Stream,
	"multiparty/from-data":     FromData,
}

// String .
func (t Type) String() string {
	if str, fond := contents[t]; fond {
		return str
	}
	return contents[DefaultContentType]
}

// Convert .
func Convert(v string) Type {
	if t, found := reverse[v]; found {
		return t
	}
	return Text
}
