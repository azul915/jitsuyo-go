package sec10

import "net/http"

func Prac() {
	http.HandleFunc("/hello", Hello)
	http.ListenAndServe(":3694", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
