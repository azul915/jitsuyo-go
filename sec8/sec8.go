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
	"github.com/ianlopshire/go-fixedwidth"
	"github.com/xuri/excelize/v2"
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

	// Unmarshal???interface{}??????????????????
	var foo interface{}
	if err := json.Unmarshal(jsonBlob, &foo); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", foo)
	// map[origin:255.255.255.255 url:https://httpbin.org/get]
	fmt.Printf("%#v\n", foo)
	// map[string]interface {}{"origin":"255.255.255.255", "url":"https://httpbin.org/get"}

	// ??????????????????
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

	in := FormInput{Name: "????????????"}
	fi, _ := json.Marshal(in)
	fmt.Println(string(fi))

	// omitempty ????????????default???(0, false, ""(?????????))??????????????????????????????Go???????????????????????????????????????
	// bottle := Bottle{
	// 	Name:        "???????????????????????????",
	// 	Price:       0,
	// 	Description: "",
	// 	HasSuger:    false,
	// }
	// out, _ := json.Marshal(bottle)
	// fmt.Println(string(out))

	bo := Bottle{
		Name:        "?????????",
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
	// Payload ??????????????????????????????????????? json.RawMessage ???????????????
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
			"message": "???????????????",
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
		{"?????????", "?????????", "????????????"},
		{"Go???????????????Web??????????????????????????????", "2016", "280"},
		{"Go???????????????????????????", "2018", "256"},
		{"Go????????????????????????????????????", "2018", "316"},
	}

	// ????????????
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

	// ????????????
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
			Name       string `csv:"??????"`
			ISOCode    string `csv:"ISO?????????"`
			Population int    `csv:"??????"`
		}
	)

	lines := []Country{
		{Name: "?????????????????????", ISOCode: "US/USA", Population: 310232863},
		{Name: "??????", ISOCode: "JP/JPN", Population: 127288000},
		{Name: "??????", ISOCode: "CN/CHN", Population: 1330044000},
	}

	// cannot use line (variable of type Country) as []string value in argument to sw.Write ???????????????
	// sw.Write????????????????????????[]string?????????

	// OpenFile?????????????????????????????????????????????????????????????????????????????????Open???Create?????????????????????
	// ????????????????????????????????????(O_RDONLY??????)???????????????????????????????????????????????????????????????????????????O_CREATE ?????????????????????????????????(umask ?????????) perm ???????????????????????????????????????????????????????????????????????????File??????????????????I/O???????????????????????????????????????????????????????????????????????? *PathError ???????????????
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
		Name       string `csv:"??????"`
		ISOCode    string `csv:"ISO?????????"`
		Population int    `csv:"??????"`
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
	s := `summary,3???
country,?????????????????????,US/USA,310232863
country,??????,JP/JPN,127288000
country,??????,CN/CHN,1330044000`

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
			fmt.Printf("summary??????????????????: %+v\n", summaries[0])
		} else {
			var countries []Country
			if err := gocsv.UnmarshalCSVWithoutHeaders(singleCSVReader{record}, &countries); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("country??????????????????: %v\n", countries[0])
		}
	}
}

func (c Country) HeaderColumns() []interface{} {
	return []interface{}{"??????", "ISO?????????", "??????"}
}

func (c Country) Columns() []interface{} {
	return []interface{}{c.Name, c.ISOCode, c.Population}
}

func Excel() {

	if _, err := os.Stat("sec8/Book1.xlsx"); err == nil {
		os.Remove("sec8/Book1.xlsx")
	}
	out := excelize.NewFile()
	out.SetCellValue("Sheet1", "A1", "Hello Excel")
	if err := out.SaveAs("sec8/Book1.xlsx"); err != nil {
		fmt.Println(err)
	}

	in, err := excelize.OpenFile("sec8/Book1.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	cell, err := in.GetCellValue("Sheet1", "A1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cell)

	lines := []Country{
		{Name: "?????????????????????", ISOCode: "US/USA", Population: 310232863},
		{Name: "??????", ISOCode: "JP/JPN", Population: 127288000},
		{Name: "??????", ISOCode: "CN/CHN", Population: 1330044000},
	}

	f := excelize.NewFile()
	sw, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range lines {
		if i == 0 {
			cell, _ := excelize.CoordinatesToCellName(1, i+1)
			sw.SetRow(cell, line.HeaderColumns())
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		sw.SetRow(cell, line.Columns())
	}

	if err := sw.Flush(); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("sec8/Book2.xlsx"); err == nil {
		os.Remove("sec8/Book2.xlsx")
	}
	if err := f.SaveAs("sec8/Book2.xlsx"); err != nil {
		log.Fatal(err)
	}

	b1, err := excelize.OpenFile("sec8/Book2.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := b1.Rows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	var countries []Country
	for i := 0; rows.Next(); i++ {
		row, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}
		if i == 0 {
			continue
		}
		population, err := strconv.Atoi(row[2])
		if err != nil {
			log.Fatal(err)
		}
		countries = append(countries, Country{
			Name:       row[0],
			ISOCode:    row[1],
			Population: population,
		})
	}
	fmt.Println(countries)

	reader, err := NewExcelCSVReader("sec8/Book2.xlsx", "Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	var cos []Country
	if err := gocsv.UnmarshalCSV(reader, &cos); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cos)
}

type excelCSVReader struct {
	rows *excelize.Rows
}

func NewExcelCSVReader(filename, sheet string) (*excelCSVReader, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	rows, err := f.Rows(sheet)
	if err != nil {
		return nil, err
	}
	return &excelCSVReader{rows}, nil
}

func (r excelCSVReader) Read() ([]string, error) {
	if r.rows.Next() {
		return r.rows.Columns()
	}
	return nil, io.EOF
}

func (r excelCSVReader) ReadAll() ([][]string, error) {
	var resp [][]string
	for r.rows.Next() {
		columns, err := r.rows.Columns()
		if err != nil {
			return nil, err
		}
		resp = append(resp, columns)
	}
	return resp, nil
}

type Book struct {
	ISBN        string
	PublishDate string
	Price       string
	PDF         string
	EPUB        string
	EbookPrice  string
}

type Book2 struct {
	ISBN        string `fixed:"1,17"`
	PublishDate string `fixed:"18,25"`
	Price       int    `fixed:"26,29"`
	PDF         string `fixed:"30,34,left"`
	EPUB        string `fixed:"35,39,left"`
	EbookPrice  int    `fixed:"40,44"`
}

func StaticLengthData() {
	s := `978-4-87311-865-9201909174620true true 3696
978-4-87311-924-3202010102750falsefalse0000
978-4-87311-878-9201903120000true true 0000`
	for _, line := range strings.Split(s, "\n") {
		r := []rune(line)

		res := Book{
			ISBN:        string(r[0:17]),
			PublishDate: string(r[17:25]),
			Price:       string(r[25:29]),
			PDF:         string(r[29:34]),
			EPUB:        string(r[34:39]),
			EbookPrice:  string(r[39:43]),
		}
		fmt.Printf("%+v\n", res)
	}

	for _, line := range strings.Split(s, "\n") {
		var b Book2
		if err := fixedwidth.Unmarshal([]byte(line), &b); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", b)
	}
}
