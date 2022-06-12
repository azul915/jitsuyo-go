package sec8

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

type Bottle struct {
	Name  string `json:"name"`
	Price int    `json:"price,omitempty"`
	KCal  *int   `json:"kcal,omitempty"`
}

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

	bottle := Bottle{
		Name:  "ミネラルウォーター",
		Price: 0,
		KCal:  Int(0),
	}
	out, _ := json.Marshal(bottle)
	fmt.Println(string(out))
}

func Int(v int) *int {
	return &v
}
