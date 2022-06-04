package main

import (
	"fmt"
	_ "jitsuyo-go/inittestb"
	_ "jitsuyo-go/inittestc"
)

func init() {
	fmt.Println("main.init")
}

func main() {
	fmt.Println("main")
}
