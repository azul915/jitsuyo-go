package sec10

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
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

// curl -X POST http://localhost:3694/comments --data '{"Message":"testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest", "UserName":"looooooooooooooongUserName"}'
// {"status":"Key: 'Comment.Message' Error:Field validation for 'Message' failed on the 'max' tag Key: 'Comment.UserName' Error:Field validation for 'UserName' failed on the 'max' tag"}
func JsonPrac() {
	type (
		Comment struct {
			Message  string `validate:"required,min=1,max=140"`
			UserName string `validate:"required,min=1,max=15"`
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
			validate := validator.New()
			if err := validate.Struct(c); err != nil {
				var out []string
				var ve validator.ValidationErrors
				if errors.As(err, &ve) {
					for _, fe := range ve {
						fmt.Printf("fe.ActualTag(): %v\n", fe.ActualTag())
						fmt.Printf("fe.Error(): %v\n", fe.Error())
						fmt.Printf("fe.Field(): %v\n", fe.Field())
						fmt.Printf("fe.Kind(): %v\n", fe.Kind())
						fmt.Printf("fe.Namespace(): %v\n", fe.Namespace())
						fmt.Printf("fe.Param(): %v\n", fe.Param())
						fmt.Printf("fe.StructField(): %v\n", fe.StructField())
						fmt.Printf("fe.StructNamespace: %v\n", fe.StructNamespace())
						fmt.Printf("fe.Tag(): %v\n", fe.Tag())
						// fmt.Printf("fe.Translate(): %v\n", fe.Translate())
						fmt.Printf("fe.Type(): %v\n", fe.Type())
						fmt.Printf("fe.Value(): %v\n", fe.Value())
						switch fe.Field() {
						case "Message":
							out = append(out, "Messageは1 ~ 140文字です")
						case "UserName":
							out = append(out, "UserNameは1 ~ 15文字です")
						}
					}
				}
				// http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusBadRequest)
				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, strings.Join(out, ",")), http.StatusBadRequest)
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
