package sec5

import (
	"errors"
	"fmt"
)

func Sec5() {
	sww := errors.New("something went wrong")
	fmt.Println(sww.Error())
}
