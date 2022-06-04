package inittestb

import (
	"fmt"
	_ "jitsuyo-go/inittesta"
)

func init() {
	fmt.Println("b.init")
}
