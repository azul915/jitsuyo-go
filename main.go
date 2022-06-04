package main

import (
	"jitsuyo-go/sec5"
	"jitsuyo-go/sec6"
)

func init() {
	sec6.Register("a", &sec6.PluginA{})
	sec6.Register("b", &sec6.PluginA{})
}

func main() {
	sec5.MultiError()
	for _, p := range sec6.Plugins() {
		p.Exec()
	}
}

// func commandLineArgs() {
// 	flag.Parse()
// 	log.Println(*sec1.FlagStr)
// 	log.Println(*sec1.FlagInt)
// 	log.Println(flag.Args())
// }
