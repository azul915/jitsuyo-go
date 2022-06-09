package sec8

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ip struct {
	Origin string `json:"origin"`
	URL    string `json:"url"`
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
}
