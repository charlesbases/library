package log

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	for level := range colors {
		fmt.Println(level.Sprint(level.short()))
	}
}
