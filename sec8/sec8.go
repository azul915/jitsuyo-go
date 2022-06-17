package sec8

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
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

func (t *JSTime) UnmarshalJSON(data []byte) error {
	var jsonNumber json.Number
	err := json.Unmarshal(data, &jsonNumber)
	if err != nil {
		return err
	}
	unix, err := jsonNumber.Int64()
	if err != nil {
		return err
	}
	*t = JSTime(time.Unix(0, unix))
	return nil
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

	s := []byte(`{
					"process_id": "001", "deleted_at": 1234567891234132
				}`)
	var rr *Record
	if err := json.Unmarshal([]byte(s), &rr); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", time.Time(rr.DeletedAt).Format(time.RFC3339Nano))

	j := `{ "time_at": "2017/01/02" }`

	var decoded JSONSample
	json.Unmarshal([]byte(j), &decoded)
	// {2017-01-02 00:00:00 +0000 UTC}
	fmt.Printf("%v\n", decoded)

	j2, _ := json.Marshal(decoded)
	// {"time_at":"2017/01/02"}
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

type Response struct {
	Type      string `json:"type"`
	Timestamp int    `json:"timestamp"`
	// Payload を具体的な構造体に展開せず json.RawMessage として保持
	Payload json.RawMessage `json:"payload"`
}

type Message struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Message   string  `json:"message"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Sensor struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Result    string `json:"result"`
	ProductID string `json:"product_id"`
}

func DecodePrac2() {
	resMessage := []byte(`{
		"type": "message",
		"timestap": 1234456,
		"payload": {
			"id": "12ct32gfds3ersadf4gfest6",
			"user_id": "ABC123",
			"message": "あいうえお",
			"latitude": 35.12321,
			"longitude": 139.12321
		}
	}`)

	resSensor := []byte(`{
		"type": "sensor",
		"timestap": 1234456,
		"payload": {
			"id": "12ct32gfds3ersadf4gfest6",
			"device_id": "ABC123",
			"result": "ok",
			"product_id": "1001"
		}
	}`)

	var rm Response
	json.Unmarshal(resMessage, &rm)

	var rs Response
	json.Unmarshal(resSensor, &rs)

	switch rm.Type {
	case "message":
		var m Message
		json.Unmarshal(rm.Payload, &m)
		fmt.Println(m)
		fmt.Println(rm)
	case "sensor":
		var s Sensor
		json.Unmarshal(rm.Payload, &s)
	}

	switch rs.Type {
	case "message":
		var m Message
		json.Unmarshal(rs.Payload, &m)
	case "sensor":
		var s Sensor
		json.Unmarshal(rs.Payload, &s)
		fmt.Println(s)
		fmt.Println(rs)
	}
}

func CSV() {
	f, err := os.Open("sec8/country.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}

	tf, _ := os.Open("sec8/country.tsv")
	defer tf.Close()

	tr := csv.NewReader(tf)
	tr.Comma = '\t'
	tr.Comment = '#'
	for {
		record, err := tr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}

	records := [][]string{
		{"書籍名", "出版年", "ページ数"},
		{"Go言語によるWebアプリケーション開発", "2016", "280"},
		{"Go言語による並列処理", "2018", "256"},
		{"Go言語でつくるインタプリタ", "2018", "316"},
	}

	// 一度消す
	if _, err := os.Stat("sec8/oreilly.csv"); err == nil {
		os.Remove("sec8/oreilly.csv")
	}

	of, err := os.OpenFile("sec8/oreilly.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()

	w := csv.NewWriter(of)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatal(err)
		}
	}
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	// 一度消す
	if _, err := os.Stat("sec8/oreilly.tsv"); err == nil {
		os.Remove("sec8/oreilly.tsv")
	}
	otf, err := os.OpenFile("sec8/oreilly.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer otf.Close()

	tw := csv.NewWriter(otf)
	tw.Comma = '\t'
	defer tw.Flush()

	for _, record := range records {
		if err := tw.Write(record); err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Error(); err != nil {
		log.Fatal(err)
	}

	type (
		Country struct {
			Name       string `csv:"国名"`
			ISOCode    string `csv:"ISOコード"`
			Population int    `csv:"人口"`
		}
	)

	lines := []Country{
		{Name: "アメリカ合衆国", ISOCode: "US/USA", Population: 310232863},
		{Name: "日本", ISOCode: "JP/JPN", Population: 127288000},
		{Name: "中国", ISOCode: "CN/CHN", Population: 1330044000},
	}

	// cannot use line (variable of type Country) as []string value in argument to sw.Write が発生する
	// sw.Writeに書き込めるのは[]stringのみ？

	// OpenFileは一般的なオープンコールで、ほとんどのユーザは代わりにOpenやCreateを使用します。
	// これは、指定されたフラグ(O_RDONLYなど)で、指定されたファイルを開く。ファイルが存在せず、O_CREATE フラグが渡された場合、(umask の前に) perm モードでファイルが作成される。成功すれば、返されたFileのメソッドをI/Oに使用することができる。エラーがある場合、それは *PathError 型である。
	// https://pkg.go.dev/os#OpenFile

	// sfo, _ := os.OpenFile("sec8/struct_country.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// sw := csv.NewWriter(sfo)
	// defer sw.Flush()
	// for _, line := range lines {
	// 	if err := sw.Write(line); err != nil {

	// 	}
	// }

	if _, err := os.Stat("sec8/struct_country.csv"); err == nil {
		os.Remove("sec8/struct_country.csv")
	}
	sf, err := os.Create("sec8/struct_country.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer sf.Close()

	if err := gocsv.MarshalFile(&lines, sf); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("sec8/struct_country_without_header.csv"); err == nil {
		os.Remove("sec8/struct_country_without_header.csv")
	}
	nf, err := os.Create("sec8/struct_country_without_header.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer nf.Close()

	if err := gocsv.MarshalWithoutHeaders(&lines, nf); err != nil {
		log.Fatal(err)
	}

	rf, err := os.Open("sec8/struct_country.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer rf.Close()

	var rlines []Country
	if err := gocsv.UnmarshalFile(rf, &rlines); err != nil {
		log.Fatal(err)
	}

	for _, v := range rlines {
		fmt.Printf("%+v\n", v)
	}

	type (
		record struct {
			Number  int    `csv:"number"`
			Message string `csv:"message"`
		}
	)

	c := make(chan interface{})
	go func() {
		defer close(c)
		for i := 0; i < 10*10; i++ {
			c <- record{i + 1, "Hello"}
		}
	}()

	if err := gocsv.MarshalChan(c, gocsv.DefaultCSVWriter(os.Stdout)); err != nil {
		log.Fatal(err)
	}

	// if _, err := os.Stat("sec8/large.csv"); err == nil {
	// 	os.Remove("sec8/large.csv")
	// }
	// lf, err := os.Open("sec8/large.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer lf.Close()

	// ch := make(chan record)
	// done := make(chan bool)
	// go func() {
	// 	if err := gocsv.UnmarshalToChan(lf, ch); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	done <- true
	// }()

	// for {
	// 	select {
	// 	case v := <-ch:
	// 		fmt.Printf("%+v\n", v)
	// 	case <-done:
	// 		return
	// 	}
	// }
}

type (
	Summary struct {
		RecordType string
		Summary    string
	}
	Country struct {
		RecordType string
		Name       string
		ISOCode    string
		Population int
	}
	singleCSVReader struct {
		record []string
	}
)

func (r singleCSVReader) Read() ([]string, error) {
	return r.record, nil
}

func (r singleCSVReader) ReadAll() ([][]string, error) {
	return [][]string{r.record}, nil
}

func MultiCSV() {
	s := `summary,3件
country,アメリカ合衆国,US/USA,310232863
country,日本,JP/JPN,127288000
country,中国,CN/CHN,1330044000`

	r := csv.NewReader(strings.NewReader(s))
	r.FieldsPerRecord = -1
	all, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range all {
		if record[0] == "summary" {
			var summaries []Summary
			if err := gocsv.UnmarshalCSVWithoutHeaders(singleCSVReader{record}, &summaries); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("summary行の読み込み: %+v\n", summaries[0])
		} else {
			var countries []Country
			if err := gocsv.UnmarshalCSVWithoutHeaders(singleCSVReader{record}, &countries); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("country行の読み込み: %v\n", countries[0])
		}
	}
}
