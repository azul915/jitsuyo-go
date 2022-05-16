package main

import (
	"flag"
	"jitsuyo-go/sec1"
	"log"
)

func main() {
	// sec1.Practice()
	// sec1.OptionArgs()
	// sec1.FunctionalOption()
	// commandLineArgs()
	sec1.PracticeSlice()
}

func commandLineArgs() {
	flag.Parse()
	log.Println(*sec1.FlagStr)
	log.Println(*sec1.FlagInt)
	log.Println(flag.Args())
}
