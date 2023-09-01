package sonyflake

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
	var count = 10000

	// uuid
	{
		start := time.Now()
		for i := 0; i < count; i++ {
			uuid.New().ID()
		}
		fmt.Println("uuid:", time.Since(start)) // 1.0408ms
	}

	// sonyflake
	{
		start := time.Now()
		for i := 0; i < count; i++ {
			NextID()
		}
		fmt.Println("sonyflake:", time.Since(start)) // 391.886ms
	}
}

func TestNextID(t *testing.T) {
	fmt.Println(uuid.New().ID())
	fmt.Println(NextID())
}
