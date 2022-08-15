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

	// バッファなしintチャネル作成
	// ic := make(chan int)

	// バッファありstringチャネル作成
	// sc := make(chan string, 10)

	// 送信
	// ic <- 0
	// sc <- "test"

	// 受信
	// 送信結果は捨てる
	// <-sc

	// 送信結果を変数に入れる
	// r := <-sc

	// 送信結果とチャネルの状態を変数に入れる
	// r, ok := <-sc

	ic := make(chan int)
	go func() {
		ic <- 10
		ic <- 20
		close(ic)
	}()

	for v := range ic {
		fmt.Println(v)
	}
}
