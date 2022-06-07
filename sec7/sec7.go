package sec7

import (
	"fmt"
	"runtime/debug"
)

func Sec7() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	// Go1.18 rc.1でfunc (*debug.BuildInfo).MarshalText() は廃止された
	str := info.String()
	fmt.Println(str)
}
