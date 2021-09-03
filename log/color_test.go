package log

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	fmt.Println(Black.String("Black"))
	fmt.Println(Red.String("Red"))
	fmt.Println(Green.String("Green"))
	fmt.Println(Yellow.String("Yellow"))
	fmt.Println(Blue.String("Blue"))
	fmt.Println(Purple.String("Purple"))
	fmt.Println(Cyan.String("Cyan"))
	fmt.Println(White.String("White"))
}
