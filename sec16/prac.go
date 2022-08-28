package sec16

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

func Prac() {
	// items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	// for _, v := range items {
	// 	go func(v int) {
	// 		fmt.Printf("v = %d\n", v)
	// 	}(v)
	// }
	// time.Sleep(time.Second)

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

// 受信型
// func recv(r <-chan string) {
// 	v := <-r
// 	r <- "送信"// 送信はNG
// }

// 送信型
// func send(s chan<- string) {
// 	s <-"送信OK"
// 	v := <-s // 受信ダメ
// }

type Account struct {
	balance int
	lock    sync.RWMutex
}

func (a *Account) GetBalance() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.balance
}

func (a *Account) Transfer(amount int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.balance += amount
}

type Account2 struct {
	balance  int
	transfer chan int
}

func NewAccount() *Account2 {
	transfer := make(chan int)
	r := &Account2{
		transfer: transfer,
	}
	go func() {
		for {
			amount := <-transfer
			r.balance += amount
		}
	}()
	return r
}

func (a *Account2) GetBalance() int {
	return a.balance
}

func (a *Account2) Transfer(amount int) {
	a.transfer <- amount
}

type Account3 struct {
	balance int64
}

func (a *Account3) GetBalance() int64 {
	return a.balance
}

func (a *Account3) Transfer(amount int64) {
	atomic.AddInt64(&a.balance, amount)
}

// type (
// 	Task   struct{}
// 	Result struct {
// 		Value interface{}
// 		Err   error
// 	}
// )

// func Worker(tasks <-chan Task, results chan<- Result) {
// 	for t := range tasks {
// 		results <- Result{
// 			Value: t,
// 			Err:   fmt.Errorf("something went wrong"),
// 		}
// 	}
// }

type (
	Task   string
	Result struct {
		Value int64
		Task  Task
		Err   error
	}
)

func worker(id int, tasks <-chan Task, results chan<- Result) {

	for t := range tasks {
		fmt.Printf("worker: %d task: %s\n", id, t)
		s, err := os.Stat(string(t))
		if err != nil && s.IsDir() {
			err = fmt.Errorf("worker %d err: %s is dir", id, string(t))
		}
		result := Result{
			Task: t,
		}
		if err != nil {
			result.Err = err
		} else {
			fmt.Printf("worker: %d path: %s size: %d\n", id, string(t), s.Size())
			result.Value = s.Size()
		}
		results <- result
	}
}

func TotalFileSize() int64 {
	tasks := make(chan Task)
	results := make(chan Result)
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(i, tasks, results)
	}
	inputDone := make(chan struct{})
	var remainedCount int64
	go func() {
		filepath.Walk(runtime.GOROOT(), func(path string, info os.FileInfo, err error) error {
			atomic.AddInt64(&remainedCount, 1)
			tasks <- Task(path)
			return nil
		})
		close(inputDone)
		close(tasks)
	}()
	var size int64

	for {
		select {
		case result := <-results:
			if result.Err != nil {
				fmt.Printf("err %v for %sl\n", result.Err, result.Task)

			} else {
				atomic.AddInt64(&size, result.Value)
			}
			atomic.AddInt64(&remainedCount, -1)
		case <-inputDone:
			if remainedCount == 0 {
				return size
			}
		}
	}
}
