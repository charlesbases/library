package base

import (
	"fmt"
	"testing"
	"time"
)

func TestBin(t *testing.T) {
	bin := NewBinary([]byte{200, 210, 220, 230, 240, 250})
	fmt.Println("十六进制:", bin.Hex())

	var loop int = 1e6

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			bin.Hex()
		}
		fmt.Println("1 >>>>:", time.Since(start)) // 271.329953ms
	}
}

func TestHex(t *testing.T) {
	hex := NewHex("0xc8d2dce6f0fa")
	fmt.Println("二进制:", hex.Binary())

	var loop int = 1e6

	{
		var start = time.Now()
		for i := 0; i < loop; i++ {
			hex.Binary()
		}
		fmt.Println("1 >>>>:", time.Since(start)) // 271.329953ms
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
