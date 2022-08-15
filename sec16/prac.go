package sec16

import (
	"fmt"
	"time"
)

func Prac() {
	fmt.Println("ゴルーチンを実行します")
	go func() {
		fmt.Println("ゴルーチンが実行されています")
	}()

	fmt.Println("ゴルーチンの終了を待ちます")
	time.Sleep(time.Second)
	fmt.Println("ゴルーチンが終了しました")
}
