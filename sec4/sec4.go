package sec4

import (
	"os"

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
