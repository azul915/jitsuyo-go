package main

import (
	"fmt"
)

var version string

// CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -X main.version=1.0.0' main.go
// ls -la ./main | awk '{ print $5 }'
// 1350992

// CGO_ENABLED=0 go build -ldflags '-X main.version=1.0.0' main.go
// ls -la ./main | awk '{ print $5 }'
// 1867088
func main() {
	fmt.Printf("Hello 世界, version = %s", version)
}
