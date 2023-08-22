package system

import (
	"fmt"
	"strconv"
)

type Number interface {
	// ToBin 二进制
	ToBin() string
	// ToDec 十进制
	ToDec() uint64
	// ToHex 十六进制
	ToHex() string
}

type number uint64

func (n number) ToBin() string {
	return fmt.Sprintf(`%b`, n)
}

func (n number) ToDec() uint64 {
	return uint64(n)
}

func (n number) ToHex() string {
	return fmt.Sprintf(`%X`, n)
}

// NewDecimal .
func NewDecimal(dec uint64) Number {
	return number(dec)
}

// NewBinary .
func NewBinary(bin string) (Number, error) {
	dec, err := strconv.ParseUint(bin, 2, 64)
	return number(dec), err
}

// NewHexadecimal .
func NewHexadecimal(hex string) (Number, error) {
	dec, err := strconv.ParseUint(hex, 16, 64)
	return number(dec), err
}
