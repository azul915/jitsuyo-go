package sec5

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/multierr"
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

// エラーハンドリング基本方針
// 1. 呼び出し元に関数の引数などの情報を付与してエラーを返す
// 2. ログを出力して処理を継続する
// 3. リトライを実施する
// 4. リソースをクローズする

type User struct{}

func getInvitedUserWithEmail(ctx context.Context, email string) (User, error) {
	return User{}, errors.New("Not Found.")
}

func ErrorHandling() error {

	// 1. 呼び出し元に関数の引数などの情報を付与してエラーを返す
	c := context.Background()
	address := "hoge@gmail.com"
	_, err := getInvitedUserWithEmail(c, address)
	if err != nil {
		return fmt.Errorf("fail to get invited user with email (%s): %w", address, err)
	}
	return nil
}

// func fetchCapacity(ctx context.Context, key string) (int, error) {
// 	var capacity int
// 	// 2. ログを出力して処理を継続する
// 	db, _ := sql.Open("driver-name", "database=test1")
// 	query := `SELECT value FROM parameter_master WHERE key = $1;`
// 	err := db.QueryRowContext(context.Background(), query, key)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {

// 		}
// 	}
// }

func MultiError() {
	var merr error
	ers := []error{nil, errors.New("Error 1"), nil, errors.New("Error 3")}
	for _, e := range ers {
		merr = multierr.Append(merr, e)
		multierr.Errors(merr)
	}

	if merr != nil {
		fmt.Println(merr)
		fmt.Println(merr.Error())
	}
}
