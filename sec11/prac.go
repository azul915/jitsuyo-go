package sec11

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func Prac() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("something went wrong!!")
	}
	resp.Body.Close()

	type User struct {
		Name string
		Addr string
	}

	u := User{
		Name: "O'Reilly Japan",
		Addr: "東京都新宿区四谷坂町",
	}
	payload, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Post("http://example.com/", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	cl := &http.Client{
		Timeout:   10 * time.Second,
		Transport: http.DefaultTransport,
	}

	nr, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	r, err := cl.Do(nr)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()
}
