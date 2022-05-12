package sec1

import (
	"errors"
)

var (
	ErrTooLong         = errors.New("bufio.Scanner: token too long")
	ErrNegativeAdvance = errors.New("bufio.Scanner: SplitFunc returns negative advance count")
	ErrAdvancetooFar   = errors.New("bufio.Scanner: SplitFunc returns advance count beyond input")
)
