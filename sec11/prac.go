package sec11

import (
	"log"
	"net/http"
)

func Prac() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("something went wrong!!")
	}
}
