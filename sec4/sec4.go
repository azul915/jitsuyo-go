package sec4

import (
	"bufio"
	"io"
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
