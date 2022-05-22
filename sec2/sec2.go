package sec2

import (
	"container/list"
	"fmt"
	"net/url"
	"time"
)

type HTTPStatus int

const (
	StatusOK              HTTPStatus = 200
	StatusUnauthorized    HTTPStatus = 401
	StatusPaymentRequired HTTPStatus = 402
	StatusForbidden       HTTPStatus = 403
)

func (s HTTPStatus) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusPaymentRequired:
		return "PaymentRequired"
	case StatusForbidden:
		return "Forbidden"
	default:
		return fmt.Sprintf("HTTPStatus(%d)", s)
	}
}

func UrlValues() {
	v1 := url.Values{}
	_ = make(url.Values)
	// url.Valuesが map[string]string なので、Add
	v1.Add("key1", "value1")
	v1.Add("key2", "value2")
	for k, v := range v1 {
		fmt.Printf("%s: %v\n", k, v)
	}
}

func ContainerList() {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	for ele := l.Front(); ele != nil; ele = ele.Next() {
		fmt.Println(ele.Value)
	}
}

type Consumer struct {
	Name       string
	ActiveFlg  bool
	ExpireDate time.Time
}

type Consumers []Consumer

func (c Consumers) ActiveConsumer() Consumers {
	resp := make(Consumers, 0, len(c))
	for _, v := range c {
		if v.ActiveFlg {
			resp = append(resp, v)
		}
	}
	return resp
}

func (c Consumers) RequiredFollow() Consumers {
	return c.ActiveConsumer().expires(time.Now().AddDate(0, 1, 0))
}

func (c Consumers) expires(end time.Time) Consumers {
	resp := make(Consumers, 0, len(c))
	for _, v := range c {
		if v.ExpireDate.Before(end) {
			resp = append(resp, v)
		}
	}
	return resp
}

func ConsumerPrac() {
	cs := Consumers{
		{
			Name:      "Tom",
			ActiveFlg: true,
		},
		{
			Name:      "Nancy",
			ActiveFlg: false,
		},
		{
			Name:      "James",
			ActiveFlg: true,
		},
	}
	acs := cs.ActiveConsumer()
	fmt.Println(acs)
}

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
