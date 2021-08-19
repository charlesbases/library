package magnifier

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	New(WithMultiple(1e10))

	var source = 1.234567899999999999999
	transformed := Magnify(source)
	fmt.Println(transformed)

	fmt.Println(Restore(transformed))
}
