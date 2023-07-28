package sonyflake

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

// 使用随机数作为 MachineID
var randMachineID = func() (uint16, error) {
	return uint16(rand.Uint32()), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())

	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: randMachineID,
	})
}

type ID uint64

// String .
func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

// NextID .
func NextID() ID {
	id, _ := sf.NextID()
	return ID(id)
}

// ParseString .
func ParseString(v string) ID {
	id, _ := strconv.ParseUint(v, 10, 16)
	return ID(id)
}
