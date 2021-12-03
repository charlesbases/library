package base

type binary bytes

// NewBinary .
func NewBinary(data []byte) *binary {
	return (*binary)(&data)
}

// Hex .
func (b *binary) Hex() string {
	var length = b.len()
	var res bytes = make([]byte, 2, length<<1+2)

	res[0], res[1] = '0', 'x'

	for i := 0; i < length; i++ {
		res = append(res, bin2Hex[(*b)[i]>>4], bin2Hex[(*b)[i]&0xf])
	}

	return res.String()
}

// Uint64 .
// 当 b > math.MaxUint64，返回 0
func (b *binary) Uint64() uint64 {
	var length = b.len()
	if length > 8 {
		return 0
	}

	var res uint64

	var bit int
	for i := length - 1; i >= 0; i-- {
		var val = (*b)[i]

		// 1 byte = 8 bit
		for bi := 0; bi < len(bits); bi++ {
			if val&bits[bi] == bits[bi] {
				res += (1 << (bit<<3 + bi))
			}
		}

		bit++
	}
	return res
}

// len .
func (b *binary) len() int {
	return len(*b)
}
