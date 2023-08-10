package nds

import (
	"fmt"
	"testing"
)

/*
01010010110100
 0011100011000

   010001101011000101101100000 | 37063520
100010001101011000101101100000 | level 13

+ + + + + +
01010010110100 | 5300
 0011100011000 | 1816
*/

func Test(t *testing.T) {
	tid := NewTileID(116.46, 39.92, 13)
	fmt.Println(fmt.Sprintf("TileID: %b", tid))
	fmt.Println("+ + + + + +")

	x, y := tid.split()
	fmt.Println(fmt.Sprintf("LEV:    %b", 1<<(16+tid.level())))
	fmt.Println(fmt.Sprintf("X:      %b", x))
	fmt.Println(fmt.Sprintf("Y:      %b", y))
	fmt.Println("+ + + + + +")

	fmt.Println("拆分再合并:", tid == tid.merge(x, y))

	{
		tid1 := tid.merge(x+1, y+1)

		x1, y1 := tid1.split()

		fmt.Println("加一再减一:", tid == tid1.merge(x1-1, y1-1))
	}
}
