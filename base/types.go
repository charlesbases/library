package base

import "unsafe"

type bytes []byte

// bits 8 bit
var bits = [8]uint8{
	0x1 /* 00000001 */, 0x1 << 1 /* 00000010 */, 0x1 << 2 /* 00000100 */, 0x1 << 3, /* 00001000 */
	0x1 << 4 /* 00010000 */, 0x1 << 5 /* 00100000 */, 0x1 << 6 /* 01000000 */, 0x1 << 7, /* 10000000 */
}

// bin2Hex binary to hexadecimal
var bin2Hex = map[uint8]byte{
	0x0: '0', 0x1: '1', 0x2: '2', 0x3: '3', 0x4: '4', 0x5: '5', 0x6: '6', 0x7: '7',
	0x8: '8', 0x9: '9', 0xa: 'a', 0xb: 'b', 0xc: 'c', 0xd: 'd', 0xe: 'e', 0xf: 'f',
}

// hex2Bin hexadecimal to binary
var hex2Bin = map[byte]uint8{
	'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3, '4': 0x4, '5': 0x5, '6': 0x6, '7': 0x7,
	'8': 0x8, '9': 0x9, 'a': 0xa, 'b': 0xb, 'c': 0xc, 'd': 0xd, 'e': 0xe, 'f': 0xf,
	'A': 0xa, 'B': 0xb, 'C': 0xc, 'D': 0xd, 'E': 0xe, 'F': 0xf,
}

func (bs *bytes) String() string {
	return *bytes2String(bs)
}

// reverse 反转、逆序
func (b *bytes) reverse() {
	var length = b.len()
	for i := 0; i < length>>1; i++ {
		(*b)[i] = (*b)[i] ^ (*b)[length-i-1]
		(*b)[length-i-1] = (*b)[i] ^ (*b)[length-i-1]
		(*b)[i] = (*b)[i] ^ (*b)[length-i-1]
	}
}

// len .
func (b *bytes) len() int {
	if b == nil {
		return 0
	}
	return len(*b)
}

func bytes2String(data *bytes) *string {
	return (*string)(unsafe.Pointer(data))
}

func string2Bytes(data *string) *bytes {
	strPointer := (*[2]uintptr)(unsafe.Pointer(data))
	bytesPointer := [3]uintptr{strPointer[0], strPointer[1], strPointer[1]}
	return (*bytes)(unsafe.Pointer(&bytesPointer))
}
