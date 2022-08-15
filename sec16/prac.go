package sec16

import (
	"fmt"
	"time"
)

func Prac() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for _, v := range items {
		v2 := v
		go func() {
			fmt.Printf("v2 = %d, address = %p\n", v2, &v2)
		}()
	}
	time.Sleep(time.Second)
}
