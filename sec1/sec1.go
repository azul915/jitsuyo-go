package sec1

import (
	"errors"
	"fmt"
)

var (
	ErrTooLong         = errors.New("bufio.Scanner: token too long")
	ErrNegativeAdvance = errors.New("bufio.Scanner: SplitFunc returns negative advance count")
	ErrAdvancetooFar   = errors.New("bufio.Scanner: SplitFunc returns advance count beyond input")
)

type (
	CarType int
)

//go:generate stringer -type=CarType

const (
	Sedan CarType = iota + 1
	Hatchback
	MPV
	SUV
	Crossover
	Coupe
	Convertible
)

const (
	E = iota + 10
	I = iota
)

const (
	a = iota // 0
	b        // 1
	c        // 2
	_        // 3だけど使われない
	// 空行は無視
	d        // 4
	e = iota // 5
)

const (
	f = iota
	g // 1
	h // 2
)

type (
	CarOption uint64
)

//go:generate stringer -type=CarOption

const (
	GPS          CarOption = 1 << iota // 1
	AWD                                // 2
	SunRoof                            // 4
	HeatedSeat                         // 8
	DriverAssist                       // 16
)

type (
	Lang int
)

const (
	Go = iota + 1
	Python
	Kotlin
	Java
	Rust
)

//go:generate enumer -type=Lang -json

type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
)

//go:generate stringer -type=Pill

func Practice() {
	fmt.Println(ErrTooLong)
	fmt.Println(ErrAdvancetooFar)
	fmt.Println(ErrNegativeAdvance)
	fmt.Println(E)
	fmt.Println(I)

	var t CarType = SUV
	fmt.Println(t)

	var o CarOption = SunRoof | HeatedSeat
	fmt.Printf("o&AWD: %v\n", o&AWD)
	fmt.Printf("o&DriverAssist: %v\n", o&DriverAssist)
	fmt.Printf("o&HeatedSeat: %v\n", o&HeatedSeat)
	fmt.Printf("o&SunRoof: %v\n", o&SunRoof)
	if o&SunRoof != 0 {
		fmt.Println("サンルーフ付き")
	}

	c := Convertible
	fmt.Printf("愛車は%sです\n", c)

	fmt.Printf("%d\n", Placebo)
	fmt.Println(Placebo)

	se := errors.New("something error")
	fmt.Println(se.Error())

	// func New(text string) error {
	// 	return &errorString{text}
	// }
	// type errorString struct {
	// 	s string
	// }
	// func (e *errorString) Error() string {
	// 	return e.s
	// }

}
