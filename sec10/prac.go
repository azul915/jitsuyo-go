package sec10

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func Prac() {
	// http.HandleFunc("/hello", Hello)
	// 	http.HandleFunc("/hello", http.HandlerFunc(Hello))
	http.Handle("/hello", HelloStruct{})
	http.ListenAndServe(":3694", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

type HelloStruct struct{}

func (h HelloStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world!")
}

// curl -X POST http://localhost:3694/comments
// {"status":"EOF"}

// curl -X GET http://localhost:3694/comments
// []

// curl -X PATCH http://localhost:3694/comments
// {"status":"permits only GET or POST"}

func JsonPrac() {
	type (
		Comment struct {
			Message  string
			UserName string
		}
	)

	var mutex = &sync.RWMutex{}
	comments := make([]Comment, 0, 100)

	http.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			mutex.RLock()

			if err := json.NewEncoder(w).Encode(comments); err != nil {
				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
				return
			}
			mutex.RUnlock()
		case http.MethodPost:
			var c Comment
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
				return
			}
			mutex.Lock()
			comments = append(comments, c)
			mutex.Unlock()

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"status":"created"}`))
		default:
			http.Error(w, `{"status":"permits only GET or POST"}`, http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":3694", nil)
}
