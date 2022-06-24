package main

import (
	"fmt"
)

type Test interface {
	A() string
	B() int
}

type Foo struct{}

func (f *Foo) A() string {
	return "test"
}
func (f *Foo) B() int {
	return 0
}

func NewTest() Test {
	return &Foo{}
}

func main() {
	// ctx := context.Background()
	// sec9.FetchUser(ctx, "")
	prac()
}

func prac() {
	// var t Test
	// t = &Foo{}
	// fmt.Println(t.A())
	// fmt.Println(t.B())

	t := NewTest()
	fmt.Println(t.A())
	fmt.Println(t.B())
}
