package sec5

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Sec5() {
	sww := errors.New("something went wrong")
	fmt.Println(sww.Error())
	ef := fmt.Errorf("something went wrong: %s", sww.Error())
	fmt.Println(ef)
}

type HTTPError struct {
	StatusCode int
	URL        string
}

// type HTTPErrorをレシーバとするError() stringを実装することでErrorインターフェースを満たす
func (he *HTTPError) Error() string {
	return fmt.Sprintf("http status code = %d, url = %s", he.StatusCode, he.URL)
}

func ReadContents(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{StatusCode: resp.StatusCode, URL: url}
	}
	return io.ReadAll(resp.Body)
}
