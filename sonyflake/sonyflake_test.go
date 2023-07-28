package sonyflake

import (
	"fmt"
	"testing"
)

func TestNextID(t *testing.T) {
	fmt.Println(NextID())
}
