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
