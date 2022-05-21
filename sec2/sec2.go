package sec2

import (
	"container/list"
	"fmt"
	"net/url"
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
