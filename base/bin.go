package base

type binary bytes

// NewBinary .
func NewBinary(data []byte) *binary {
	return new(binary).handle(data)
}

// handle
func (b *binary) handle(data []byte) *binary {
	// 去除首部无效的 0
	if data != nil {
		for idx, val := range data {
			if val != 0 {
				*b = data[idx:]
				return b
			}
		}
	}

	*b = []byte{}
	return b
}

// Decimal binary to decimal
// 当 b > math.MaxUint64，返回 0
func (b *binary) Decimal() uint64 {
	if b.IsNull() || b.len() > 8 {
		return 0
	}

	var res uint64

	// []byte 从右到左位置
	var bit int

	var length = b.len()
	for i := length - 1; i >= 0; i-- {
		var val = (*b)[i]

		// 1 byte = 8 bit
		for bi := 0; bi < len(bits); bi++ {
			if val&bits[bi] == bits[bi] {
				res += 1 << (bit<<3 + bi)
			}
		}

		bit++
	}
	return res
}

// Hexadecimal binary to hexadecimal
func (b *binary) Hexadecimal() string {
	if b.IsNull() {
		return ""
	}

	var length = b.len()

	var res bytes = make([]byte, 2, length<<1+2)
	res[0], res[1] = '0', 'x'

	for i := 0; i < length; i++ {
		res = append(res, bin2Hex[(*b)[i]>>4], bin2Hex[(*b)[i]&0xf])
	}

	return res.String()
}

// Display .
func (b *binary) Display() string {
	if b.IsNull() {
		return ""
	}

	var res bytes
	var length = b.len()

	res = make([]byte, 0, length<<1)

	for i := 0; i < length; i++ {
		var val = (*b)[i]
		for bi := len(bits) - 1; bi >= 0; bi-- {
			if val&bits[bi] == bits[bi] {
				res = append(res, '1')
			} else {
				res = append(res, '0')
			}
		}
	}

	return res.String()
}

// IsNull .
func (b *binary) IsNull() bool {
	return b == nil || b.len() == 0
}

// len .
func (b *binary) len() int {
	return len(*b)
}
