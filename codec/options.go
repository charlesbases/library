package codec

// MarshalOptions .
type MarshalOptions struct {
	Indent bool
}

type MarshalOption func(o *MarshalOptions)

// MarshalIndent .
func MarshalIndent() MarshalOption {
	return func(o *MarshalOptions) {
		o.Indent = true
	}
}
