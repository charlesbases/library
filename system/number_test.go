package system

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestNumber(t *testing.T) {
	var n Number

	var loop int = 1e6

	{
		fns := map[string]func(){
			"bin": func() { // 30.5032ms
				n, _ = NewBinary("1111111111111111111111111111111111111111111111111111111111111111")
			},
			"dec": func() { // 13.6812ms
				n = NewDecimal(math.MaxUint64)
			},
			"hex": func() { // 93.5948ms
				n, _ = NewHexadecimal("FFFFFFFFFFFFFFFF")
			},
		}

		for name, fn := range fns {
			start := time.Now()
			for i := 0; i < loop; i++ {
				fn()
			}
			fmt.Println(fmt.Sprintf(`new(%s): %v`, name, time.Since(start)))
		}
	}

	fmt.Println()

	{
		fns := map[string]func(){
			"bin": func() { // 139.7205ms
				n.ToBin()
			},
			"dec": func() { // 2.2277ms
				n.ToDec()
			},
			"hex": func() { // 106.5134ms
				n.ToHex()
			},
		}

		for name, fn := range fns {
			start := time.Now()
			for i := 0; i < loop; i++ {
				fn()
			}
			fmt.Println(fmt.Sprintf(`%s: %v`, name, time.Since(start)))
		}
	}
}
