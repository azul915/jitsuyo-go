package sec1

import (
	"errors"
	"flag"
	"fmt"
	"time"
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

type Portion int

const (
	Regular Portion = iota
	Small
	Large
)

type Udon struct {
	men      Portion
	aburaage bool
	ebiten   uint
}

func NewUdon(p Portion, aburaage bool, ebiten uint) *Udon {
	return &Udon{
		men:      p,
		aburaage: aburaage,
		ebiten:   ebiten,
	}
}

func NewKakeUdon(p Portion) *Udon {
	return &Udon{
		men:      p,
		aburaage: false,
		ebiten:   0,
	}
}

func NewKitsuneUdon(p Portion) *Udon {
	return &Udon{
		men:      p,
		aburaage: true,
		ebiten:   0,
	}
}

func NewTempuraUdon(p Portion) *Udon {
	return &Udon{
		men:      p,
		aburaage: false,
		ebiten:   3,
	}
}

type Option struct {
	men      Portion
	aburaage bool
	ebiten   uint
}

// dislike👎
func NewUdonWithOpt(opt Option) *Udon {
	// ゼロ値に対するデフォルト値処理は関数内部で行う

	// 朝食時間は海老天1本無料
	if opt.ebiten == 0 && time.Now().Hour() < 10 {
		opt.ebiten = 1
	}
	return &Udon{
		men:      opt.men,
		aburaage: opt.aburaage,
		ebiten:   opt.ebiten,
	}
}

type fluentOpt struct {
	men      Portion
	aburaage bool
	ebiten   uint
}

func NewUdonWithFlu(p Portion) *fluentOpt {
	// これをデフォルトとする
	return &fluentOpt{
		men:      p,
		aburaage: false,
		ebiten:   1,
	}
}

func (o *fluentOpt) Aburaage() *fluentOpt {
	o.aburaage = true
	return o
}

func (o *fluentOpt) Ebiten(n uint) *fluentOpt {
	o.ebiten = n
	return o
}

func (o *fluentOpt) Order() *Udon {
	return &Udon{
		men:      o.men,
		aburaage: o.aburaage,
		ebiten:   o.ebiten,
	}
}

func useFluentInterface() *Udon {
	return NewUdonWithFlu(Large).Aburaage().Order()
}

func OptionArgs() {
	var tempraUdon = NewUdon(Large, false, 2)
	kakeUdon := NewKakeUdon(Small)
	kitsuneUdon := NewKitsuneUdon(Regular)
	fmt.Printf("templaUdon: %v\n", tempraUdon)
	fmt.Printf("kakeUdon: %v\n", kakeUdon)
	fmt.Printf("kitsuneUdon: %v\n", kitsuneUdon)
	breakFastUdon := NewUdonWithOpt(Option{Small, false, 0})
	fmt.Printf("breakFastUdon: %v\n", breakFastUdon)
	oomoriKitsune := useFluentInterface()
	fmt.Printf("oomoriKitsune: %v\n", oomoriKitsune)
}

type OptFunc func(r *Udon)

func NewUdonFunctionaOption(opts ...OptFunc) *Udon {
	r := &Udon{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func OptMen(p Portion) OptFunc {
	return func(r *Udon) { r.men = p }
}

func OptAburaage() OptFunc {
	return func(r *Udon) { r.aburaage = true }
}

func OptEbiten(n uint) OptFunc {
	return func(r *Udon) { r.ebiten = n }
}

func FunctionalOption() {
	tokuseiUdon := NewUdonFunctionaOption(OptAburaage(), OptEbiten(3))
	fmt.Printf("tokuseiUdon: %v\n", tokuseiUdon)
}

var (
	FlagStr = flag.String("string", "default", "文字列フラグ")
	FlagInt = flag.Int("int", -1, "数値フラグ")
)
