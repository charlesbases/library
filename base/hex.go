package base

import (
	"library"
)

type hex bytes

// NewHex .
func NewHex(data string) *hex {
	if library.REGEXP_HEX.MatchString(&data) {
		var hex hex = []byte{'0', 'x'}
		if len(data) > 2 && data[0] == '0' && data[1] == 'x' {
			hex = append(hex, (*string2Bytes(&data))[2:]...)
		} else {
			hex = append(hex, *string2Bytes(&data)...)
		}
		return &hex
	}
	return new(hex)
}

// Binary hexadecimal to binary
func (h *hex) Binary() []byte {
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

func (h *hex) len() int {
	return len(*h)
}
