package rootpath

import (
	"fmt"
	"testing"
)

var root = `E:\Programs\src\github.com\charlesbases`

func TestRoot(t *testing.T) {
	r := NewRoot(root)

	{
		roots, _ := r.Files(func(o *Options) {
			o.MaxDepth = 2
		})
		for _, val := range roots {
			fmt.Println(val.String())
		}
	}
}
