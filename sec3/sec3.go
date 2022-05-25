package sec3

import (
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

type Book struct {
	Title      string
	Author     string
	Publisher  string
	ISBN       string
	ReleasedAt time.Time
}

func Dereference() {
	b := &Book{
		Title: "Mithril",
	}
	fmt.Println(b.Title)
	fmt.Println((*b).Title)

	b2 := &b
	// fmt.Println(b2.Title) //NG
	fmt.Println((**b2).Title)
	fmt.Println((*b2).Title)
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func CreateInstance() {

	// *Person型、ゼロ値
	p1 := new(Person)

	// Person型、ゼロ値
	var p2 Person

	// *Person型、設定された初期値
	p3 := &Person{
		FirstName: "Taro",
		LastName:  "Yamada",
	}

	// Person型、設定された初期値
	p4 := Person{
		FirstName: "Ichiro",
		LastName:  "Suzuki",
	}
	fmt.Printf("p1: %v, p2: %v, p3: %v, p4: %v\n", p1, p2, *p3, p4)

	// *Person型、nilが代入
	var p5 *Person
	fmt.Println(p5)
	// *Person型のp5にはnilが入っているので、
	// デリファレンスしようとするとpanic(invalid memory address or nil pointer dereference)が発生する
	// fmt.Println(*p5)
}

func (p Person) SetNameA(first, last string) {
	p.FirstName = first
	p.LastName = last
}

func (p *Person) SetNameB(first, last string) {
	p.FirstName = first
	p.LastName = last
}

func (p Person) SetNameC(first, last string) Person {
	var rtn Person
	rtn.FirstName = first
	rtn.LastName = last
	return rtn
}

func SetNamePrac() {
	var p1 Person
	p1.SetNameA("Taro", "Yamada")
	fmt.Printf("p1: %v\n", p1)

	p2 := new(Person)
	p2.SetNameB("Taro", "Yamada")
	fmt.Printf("p2: %v\n", p2)

	p3 := Person{}
	p4 := p3.SetNameC("Taro", "Yamada")
	fmt.Printf("p3: %v, p4: %v\n", p3, p4)
}

type StructWithPointer struct {
	v *int
}

func (a StructWithPointer) Modify() {
	*a.v = 10
}

// Bad😇
// インスタンスレシーバなのにインスタンスそのものに変更ができてしまう
func Foo() {
	s := StructWithPointer{}
	i := 1
	s.v = &i
	fmt.Printf("s: %v\n", *s.v)
	s.Modify()
	fmt.Printf("s: %v\n", *s.v)
}

type People []Person

func (p People) SeijinA() People {
	var resp People
	for _, e := range p {
		if 18 < e.Age {
			resp = append(resp, e)
		}
	}
	return resp
}

func Bar() {
	pp := People{
		{
			FirstName: "Taro",
			LastName:  "Yamada",
			Age:       17,
		},
		{
			FirstName: "Jiro",
			LastName:  "Tanaka",
			Age:       18,
		},
		{
			FirstName: "Hiroshi",
			LastName:  "Suzuki",
			Age:       20,
		},
	}

	fmt.Printf("pp: %v\n", pp)
	sj := pp.SeijinA()
	fmt.Printf("pp: %v\n", pp)
	fmt.Printf("sj: %v\n", sj)
}

func String[T any](s T) string {
	return fmt.Sprintf("%v", s)
}

type Struct struct {
	t interface{}
}

func (s Struct) String() string {
	return fmt.Sprintf("%v", s.t)
}

// func (s Struct[T]) Method[R any](r R) {
// 	fmt.Println(s.t, r)
// }

// 埋め込みはは名前をつけずに型だけを構造体のフィールドとして記述
type OreillyBook struct {
	Book
	ISBN13 string
}

func (b Book) GetAmazonURL() string {
	return fmt.Sprintf("https://amazon.co.jp/dp/%s", b.ISBN)
}

func (o OreillyBook) GetOreillyURL() string {
	return fmt.Sprintf("https://www.oreilly.co.jp/books/%s/", o.ISBN13)
}

func Embed() {
	ob := OreillyBook{
		ISBN13: "9784873119038",
		Book: Book{
			ISBN:  "123456789",
			Title: "Real World HTTP",
		},
	}
	fmt.Println(ob.GetAmazonURL())
	fmt.Println(ob.GetOreillyURL())

	// of course, this is also ok
	fmt.Println(ob.Book.GetAmazonURL())
}

func Hoge() {
	type (
		MapStruct struct {
			Str     string  `map:"str"`
			StrPtr  *string `map:"str"`
			Bool    bool    `map:"bool"`
			BoolPtr *bool   `map:"bool"`
			Int     int     `map:"int"`
			IntPtr  *int    `map:"int"`
		}
	)
	src := map[string]string{
		"str":  "string data",
		"bool": "true",
		"int":  "12345",
	}
	var ms MapStruct
	Decode(&ms, src)
	log.Println(ms)

	dest := map[string]string{}
	mms := MapStruct{
		Str:     "string-value",
		StrPtr:  &[]string{"string-ptr-value"}[0],
		Bool:    true,
		BoolPtr: &[]bool{true}[0],
		Int:     12345,
		IntPtr:  &[]int{12345}[0],
	}
	Encode(dest, &mms)
	log.Println(dest)
}

func Decode(target interface{}, src map[string]string) error {
	v := reflect.ValueOf(target)
	e := v.Elem()
	return decode(e, src)
}

func decode(e reflect.Value, src map[string]string) error {
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			if err := decode(e.Field(i), src); err != nil {
				return err
			}
			continue
		}
		if f.Type.Kind() == reflect.Struct {
			if err := decode(e.Field(i), src); err != nil {
				return err
			}
			continue
		}
		key := f.Tag.Get("map")
		if key == "" {
			key = f.Name
		}

		sv, ok := src[key]
		if !ok {
			continue
		}

		var k reflect.Kind
		var isP bool
		if f.Type.Kind() != reflect.Ptr {
			k = f.Type.Kind()
		} else {
			k = f.Type.Elem().Kind()
			if k == reflect.Ptr {
				continue
			}
			isP = true
		}

		switch k {
		case reflect.String:
			if isP {
				e.Field(i).Set(reflect.ValueOf(&sv))
			} else {
				e.Field(i).SetString(sv)
			}
		case reflect.Bool:
			b, err := strconv.ParseBool(sv)
			if err == nil {
				if isP {
					e.Field(i).Set(reflect.ValueOf(&b))
				} else {
					e.Field(i).SetBool(b)
				}
			}
		case reflect.Int:
			n64, err := strconv.ParseInt(sv, 10, 64)
			if err == nil {
				if isP {
					n := int(n64)
					e.Field(i).Set(reflect.ValueOf(&n))
				} else {
					e.Field(i).SetInt(n64)
				}
			}
		}
	}
	return nil
}

func Encode(target map[string]string, src interface{}) error {
	v := reflect.ValueOf(src)
	e := v.Elem()
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			if err := Encode(target, e.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}
		key := f.Tag.Get("map")
		if key == "" {
			key = f.Name
		}
		if f.Type.Kind() == reflect.Struct {
			Encode(target, e.Field(i).Addr().Interface())
			continue
		}
		var k reflect.Kind
		var isP bool
		if f.Type.Kind() != reflect.Ptr {
			k = f.Type.Kind()
		} else {
			k = f.Type.Elem().Kind()
			isP = true
			if k == reflect.Ptr {
				continue
			}
		}

		switch k {
		case reflect.String:
			if isP {
				if e.Field(i).Pointer() != 0 {
					target[key] = *(*string)(unsafe.Pointer(e.Field(i).Pointer()))
				}
			} else {
				target[key] = e.Field(i).String()
			}
		case reflect.Bool:
			var b bool
			if isP {
				if e.Field(i).Pointer() != 0 {
					b = *(*bool)(unsafe.Pointer(e.Field(i).Pointer()))
				}
			} else {
				b = e.Field(i).Bool()
			}
			target[key] = strconv.FormatBool(b)
		case reflect.Int:
			var n int64
			if isP {
				if e.Field(i).Pointer() != 0 {
					n = int64(*(*int)(unsafe.Pointer(e.Field(i).Pointer())))
				}
			} else {
				n = e.Field(i).Int()
			}
			target[key] = strconv.FormatInt(n, 10)
		}
	}
	return nil
}

type NoCopyStruct struct {
	self  *NoCopyStruct
	Value *string
}

func NewNoCopyStruc(value string) *NoCopyStruct {
	r := &NoCopyStruct{
		Value: &value,
	}
	r.self = r
	return r
}

func (n *NoCopyStruct) String() string {
	if n != n.self {
		panic("should not copy NoCopyStruct instance without Copy() method")
	}
	return *n.Value
}

func (n *NoCopyStruct) Copy() *NoCopyStruct {
	str := *n.Value
	p2 := &NoCopyStruct{
		Value: &str,
	}
	p2.self = p2
	return p2
}

func NoCopyStructPrac() {
	nnc1 := NewNoCopyStruc("practice")
	fmt.Println(nnc1.String())
	nnc2 := NoCopyStruct{}
	fmt.Println(nnc2.String())
}

type Currency struct{}
type MutableMoney struct {
	currency Currency
	amount   *big.Int
}

func (m MutableMoney) Currency() Currency {
	return m.currency
}

func (m *MutableMoney) SetCurrency(c Currency) {
	m.currency = c
}

type ImmutablMoney struct {
	currency Currency
	amount   *big.Int
}

func (im ImmutablMoney) Currency() Currency {
	return im.currency
}

func (im ImmutablMoney) SetCurrency(c Currency) ImmutablMoney {
	return ImmutablMoney{
		currency: c,
		amount:   im.amount,
	}
}

func ChanPrac() {
	wait := make(chan struct{})
	go func() {
		fmt.Println("send")
		wait <- struct{}{}
	}()
	fmt.Println("wait")
	<-wait
	fmt.Println("finished")
}

var pool *sync.Pool

type BigStruct struct {
	Member string
}

func Pool() {
	pool = &sync.Pool{
		New: func() interface{} {
			return &BigStruct{}
		},
	}

	b := pool.Get().(*BigStruct)
	pool.Put(b)
}

func NewBigStruct() *BigStruct {
	b := pool.Get().(*BigStruct)
	b.Member = ""
	return b
}

type Parent struct{}

func (p Parent) m1() {
	p.m2()
}

func (p Parent) m2() {
	fmt.Println("Parent")
}

type Child struct {
	Parent
}

func (c Child) m2() {
	fmt.Println("Child")
}

func EmbedPrac() {
	c := Child{}
	c.m1() // Parent
	c.m2() // Child
}
