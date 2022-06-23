package main

import (
	"context"
	"jitsuyo-go/sec9"
)

func main() {
	ctx := context.Background()
	sec9.FetchUser(ctx, "")
}
