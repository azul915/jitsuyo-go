package sec8

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type ip struct {
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type user struct {
	UserID    string   `json:"user_id"`
	UserName  string   `json:"user_name"`
	Languages []string `json:"languages"`
}

type FormInput struct {
	Name        string `json:"name"`
	CompanyName string `json:"company_name,omitempty"`
}
type Bottle struct {
	Name        string  `json:"name"`
	Price       int     `json:"price,omitempty"`
	KCal        *int    `json:"kcal,omitempty"`
	Description *string `json:"description,omitempty"`
	HasSuger    *bool   `json:"hasSuger,omitempty"`
}

type Rectangle struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Record struct {
	ProcessID string `json:"process_id"`
	DeletedAt JSTime `json:"deleted_at"`
}
type JSTime time.Time

type MyTime struct {
	*time.Time
}

func (mt *MyTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006/01/02\"", string(data))
	*mt = MyTime{&t}
	return err
}

func (mu MyTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(mu.Format("2006/01/02"))
}

func (t JSTime) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte("null"), nil
	}
	v := strconv.Itoa(int(tt.UnixMilli()))
	return []byte(v), nil
}

type JSONSample struct {
	TimeAt MyTime `json:"time_at"`
}

// func (t *JSTime) UnmarshalJSON(data []byte) error {
// 	var jsonNumber json.Number

// }
func DecodePrac() {
	f, err := os.Open("sec8/ip.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var resp ip
	if err := json.NewDecoder(f).Decode(&resp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)

	jsonBlob := []byte(`{
		"origin": "255.255.255.255",
		"url": "https://httpbin.org/get"
	}`)
	var hoge ip
	if err := json.Unmarshal(jsonBlob, &hoge); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", hoge)

	// Unmarshalはinterface{}も受け取れる
	var foo interface{}
	if err := json.Unmarshal(jsonBlob, &foo); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", foo)
	// map[origin:255.255.255.255 url:https://httpbin.org/get]
	fmt.Printf("%#v\n", foo)
	// map[string]interface {}{"origin":"255.255.255.255", "url":"https://httpbin.org/get"}

	// キャストする
	origin := foo.(map[string]interface{})["origin"].(string)
	fmt.Println(origin)
	url := foo.(map[string]interface{})["url"].(string)
	fmt.Println(url)

	var b bytes.Buffer
	u := user{
		UserID:   "001",
		UserName: "gopher",
	}
	json.NewEncoder(&b).Encode(u)
	fmt.Printf("%v\n", b.String())

	m, _ := json.Marshal(u)
	fmt.Println(string(m))

	uu := user{
		UserID: "001", UserName: "gopher", Languages: []string{},
	}
	bb, _ := json.Marshal(uu)
	fmt.Println(string(bb))

	in := FormInput{Name: "山田太郎"}
	fi, _ := json.Marshal(in)
	fmt.Println(string(fi))

	// omitempty は各型のdefault値(0, false, ""(空文字))を明示的に指定してもGoのゼロ値でエンコードされる
	// bottle := Bottle{
	// 	Name:        "ミネラルウォーター",
	// 	Price:       0,
	// 	Description: "",
	// 	HasSuger:    false,
	// }
	// out, _ := json.Marshal(bottle)
	// fmt.Println(string(out))

	bo := Bottle{
		Name:        "飲み物",
		Price:       0,
		KCal:        Int(0),
		Description: String(""),
		HasSuger:    Bool(false),
	}
	fmt.Println(bo.String())

	ufBlob := []byte(`{
		"width": 5,
		"height": 10,
		"radius": 6
	}`)
	var rect Rectangle
	d := json.NewDecoder(bytes.NewReader(ufBlob))
	d.DisallowUnknownFields()
	if err := d.Decode(&rect); err != nil {
		// json: unknown field "radius"
		fmt.Println(err)
	}

	r := &Record{
		ProcessID: "0001",
		DeletedAt: JSTime(time.Now()),
	}
	bbb, _ := json.Marshal(r)
	fmt.Println(string(bbb))

	j := `{ "time_at": "2017/01/02" }`

	var decoded JSONSample
	json.Unmarshal([]byte(j), &decoded)
	fmt.Printf("%v\n", decoded)

	j2, _ := json.Marshal(decoded)
	fmt.Printf("%v\n", string(j2))
}

func Int(v int) *int {
	return &v
}

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func (b *Bottle) String() string {
	return fmt.Sprintf("{ %v, %v, %v, %v, %v }", b.Name, b.Price, *b.KCal, *b.Description, *b.HasSuger)
}
