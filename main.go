package main

import (
	"fmt"
	"jitsuyo-go/sec2"
)

func main() {

	st := sec2.StatusOK
	fmt.Println(st.String())
}

// func commandLineArgs() {
// 	flag.Parse()
// 	log.Println(*sec1.FlagStr)
// 	log.Println(*sec1.FlagInt)
// 	log.Println(flag.Args())
// }
