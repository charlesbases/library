package base

import (
	"fmt"
	"testing"
	"time"
)

func TestBin(t *testing.T) {
	bin := NewBinary([]byte{255, 255, 255, 255, 255, 255, 255, 255})
	fmt.Println("二进制:", bin.Display())
	fmt.Println("十进制:", bin.Decimal())
	fmt.Println("十六进制:", bin.Hexadecimal())

	var loop int = 1e6

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			bin.Decimal()
		}
		fmt.Println("1 >>>>:", time.Since(start)) // 112.864827ms
	}

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			bin.Hexadecimal()
		}
		fmt.Println("2 >>>>:", time.Since(start)) // 271.329953ms
	}
}

func TestDec(t *testing.T) {
	var dec = NewDecimal(18446744073709551615)
	fmt.Println("二进制:", string(dec.Binary()))
	fmt.Println("十进制:", dec.Display())
	fmt.Println("十六进制:", dec.Hexadecimal())

	var loop int = 1e6

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			dec.Binary()
		}
		fmt.Println("1 >>>>:", time.Since(start)) // 112.864827ms
	}

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			dec.Hexadecimal()
		}
		fmt.Println("2 >>>>:", time.Since(start)) // 271.329953ms
	}
}

func TestHex(t *testing.T) {
	var hex = NewHexadecimal("0xffffffffffffffff")
	fmt.Println("二进制:", hex.Binary())
	fmt.Println("十进制:", hex.Decimal())
	fmt.Println("十六进制:", hex.Display())

	var loop int = 1e6

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			hex.Binary()
		}
		fmt.Println("1 >>>>:", time.Since(start)) // 112.864827ms
	}

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			hex.Decimal()
		}
		fmt.Println("2 >>>>:", time.Since(start)) // 271.329953ms
	}
}

func TestDemo(t *testing.T) {
	var loop int = 1e6

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {

		}
		fmt.Println("1 >>>>:", time.Since(start))
	}

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {

		}
		fmt.Println("2 >>>>:", time.Since(start))
	}
}
