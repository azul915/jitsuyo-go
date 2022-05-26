package sec4

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gen2brain/beeep"
)

type Warning interface {
	Show(message string)
}

type DesktopWarning struct{}

func (d DesktopWarning) Show(message string) {
	beeep.Alert(os.Args[0], message, "")
}

func NewDesktopWarning() Warning {
	return DesktopWarning{}
}

func IntefacePrac() {
	// cannot use &(DesktopWarning literal) (value of type *DesktopWarning) as Warning value in variable declaration: *DesktopWarning does not implement Warning (missing method Show)c
	// var _ Warning = &DesktopWarning{}

	warn := DesktopWarning{}
	warn.Show("Hello World to desktop")

	wn := NewDesktopWarning()
	wn.Show("Hello World to desktop")
}

func Normalize(w io.Writer, r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		// s, err := br.ReadString('\n')
		s, _, err := br.ReadLine()
		if string(s) != "" {
			io.WriteString(w, string(s))
		}
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func NormalizeFile(input, output string) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	return Normalize(w, r)
}

func NormalizeString(i string) (string, error) {
	r := strings.NewReader(i)
	// var w strings.Builder
	w := new(strings.Builder)
	err := Normalize(w, r)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func CastPrac() {
	ctx := context.WithValue(context.Background(), "favorite", "銭形平次")

	// ctx.Value()はinterface{}
	// okでキャスト成功可否を確認する
	if s, ok := ctx.Value("favorite").(string); ok {
		log.Printf("私が好きなものは%sです\n", s)
	}

	switch v := ctx.Value("favorite").(type) {
	case string:
		log.Printf("好きなものは: %s\n", v)
	case int:
		log.Printf("好きな数値は: %d\n", v)
	case complex128:
		log.Printf("好きな複素数は: %f\n", v)
	default:
		log.Printf("好きなものは: %v\n", v)
	}
}

func FishListPrac() {

	// type Fish interface{}
	// var fishList = []Fish{"鯖", "鰤", "鮪"}
	// var fishNameList = fishList.([]string) // Compile Error
	// var anyList []any = fishList // Compile Error

	var fishList = []any{"鯖", "鰤", "鮪"}
	fishNames := make([]string, len(fishList))
	for i, f := range fishList {
		if fn, ok := f.(string); ok {
			fishNames[i] = fn
		}
	}

	fibonacciNumbers := []int{1, 1, 2, 3, 5, 8}
	anyValues := make([]any, len(fibonacciNumbers))
	for i, fn := range fibonacciNumbers {
		// アップキャスト型アサーション不要
		anyValues[i] = fn
	}
}
