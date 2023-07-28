package base

type decimal uint64

// NewDecimal .
func NewDecimal(dec uint64) *decimal {
	return (*decimal)(&dec)
}

// Binary decimal to binary
func (d *decimal) Binary() []byte {
	if d.IsZero() {
		return []byte{0}
	}

	var bin bytes = make([]byte, 0, 8)
	var num = uint64(*d)
	for ; num > 0; num >>= 1 {
		bin = append(bin, dec2Bin[num%2])
	}
	bin.reverse()
	return bin
}

// Hexadecimal decimal to hexadecimal
func (d *decimal) Hexadecimal() string {
	if d.IsZero() {
		return "0x0"
	}

	var hex bytes = make([]byte, 0, 8)
	var num = uint64(*d)
	for ; num > 0; num >>= 4 {
		hex = append(hex, dec2Hex[num%16])
	}
	hex = append(hex, 'x', '0')
	hex.reverse()
	return hex.String()
}

// Display .
func (d *decimal) Display() uint64 {
	return uint64(*d)
}

// IsZero .
func (d *decimal) IsZero() bool {
	return d == nil || *d == 0
}
