package sec10

import "net/http"

func Prac() {
	http.HandlerFunc(Hello)

}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
