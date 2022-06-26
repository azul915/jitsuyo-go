package sec10

import "net/http"

func Prac() {
	// http.HandleFunc("/hello", Hello)
	http.Handle("/hello", HelloStruct{})
	http.ListenAndServe(":3694", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

type HelloStruct struct{}

func (h HelloStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
