package xerror

import (
	"fmt"
	"testing"
)

// ecode .
func ecode() error {
	return ServerErr
}

// xerror .
func xerror() error {
	return New(ServerErr, "数据库连接失败")
}

func Test(t *testing.T) {
	fmt.Println(ecode())
	fmt.Println(xerror())
}
