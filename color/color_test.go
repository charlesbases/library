package color

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	color := New(FgCyan, BgGreen)
	fmt.Println(color.Sprintf("test"))
}
