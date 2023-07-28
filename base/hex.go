package base

import "github.com/charlesbases/library/regexp"

type hexadecimal bytes

// NewHexadecimal .
func NewHexadecimal(data string) *hexadecimal {
	switch regexp.HEX.MatchString(data) {
	case true:
		var hex hexadecimal = []byte{'0', 'x'}
		if len(data) > 2 && data[0] == '0' && data[1] == 'x' {
			hex = append(hex, (*string2Bytes(&data))[2:]...)
		} else {
			hex = append(hex, *string2Bytes(&data)...)
		}
		return &hex
	default:
		return &hexadecimal{}
	}
}

// Binary hexadecimal to binary
func (h *hexadecimal) Binary() []byte {
	if h.IsNull() {
		return []byte{0}
	}

	var bin bytes
	var length = h.len()
	switch length {
	// 0x?
	case 3:
		bin = append(bin, hex2Bin[(*h)[2]])
	default:
		var start = 2
		if length%2 == 1 {
			bin = append(bin, hex2Bin[(*h)[2]])
			start = 3
		}

		for i := start; i < length; i += 2 {
			bin = append(bin, (hex2Bin[(*h)[i]]<<4)|hex2Bin[(*h)[i+1]])
		}
	}

	return bin
}

// Decimal hexadecimal to decimal
// hexadecimal 大于 math.MaxUint64 时，会产生溢出
func (h *hexadecimal) Decimal() uint64 {
	var hex uint64
	var bit int
	for i := h.len() - 1; i > 1; i-- {
		hex += hex2Dec[(*h)[i]] * pow(16, bit)
		bit++
	}
	return hex
}

// Display .
func (h *hexadecimal) Display() string {
	return (*bytes)(h).String()
}

// IsNull .
func (b *hexadecimal) IsNull() bool {
	return b == nil || b.len() == 0
}

// len .
func (h *hexadecimal) len() int {
	return len(*h)
}

// pow return x^y
func pow(x, y int) uint64 {
	switch {
	case x == 0:
		return 0
	case y == 0, x == 1:
		return 1
	case y == 1:
		return uint64(x)
	default:
		var res uint64 = 1
		for i := 1; i <= y; i++ {
			res *= uint64(x)
		}
		return res
	}
}
