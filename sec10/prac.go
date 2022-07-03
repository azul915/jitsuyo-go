package sec10

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func MustCheck() {
	// 必須チェックを行いたいフィールドをポインタ型で定義
	type Book struct {
		Title string `validate:"required"`
		Price *int   `validate:"required"`
	}

	s := `{"Title":"Real World HTTP ミニ版", "Price": 0}`
	var b Book
	if err := json.Unmarshal([]byte(s), &b); err != nil {
		log.Fatal(err)
	}

	if err := validator.New().Struct(b); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fmt.Printf("フィールド %s が %s 違反です (値: %v) \n", fe.Field(), fe.Tag(), fe.Value())
			}
		}
	}
}

// curl -X GET -v -G -d "searchword=検索用語" -d other=value http://localhost:3694/params
func Param() {
	http.HandleFunc("/params", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		word := r.FormValue("searchword")
		log.Printf("searchword: %s\n", word)

		// _, ok = map[key]でkeyがmapに含まれるかどうかをboolで取れる
		words, ok := r.Form["searchword"]
		log.Printf("search word = %v has values %v\n", words, ok)

		log.Print("all queries")
		for k, v := range r.Form {
			log.Printf("   %s: %s\n", k, v)
		}
	})

	// touch index.rst && curl -F file=@index.rst -F data=other http://localhost:3694/file && rm index.rst
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 * 1024 * 1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f, h, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(h.Filename)
		o, err := os.Create(h.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer o.Close()
		_, err = io.Copy(o, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		value := r.PostFormValue("data")
		log.Printf(" value = %s", value)
	})
	http.ListenAndServe(":3694", nil)
}
