package sec16

import (
	"fmt"
	"time"
)

func Prac() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for _, v := range items {
		go func(v int) {
			fmt.Printf("v = %d\n", v)
		}(v)
	}
	time.Sleep(time.Second)
}
