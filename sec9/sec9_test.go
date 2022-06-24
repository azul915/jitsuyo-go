package sec9

import (
	"context"
	"database/sql"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestFetchUser(t *testing.T) {
	connStr := "host=localhost port=5432 user=user dbname=db password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlBytes, err := ioutil.ReadFile("./sec9/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(context.TODO(), string(sqlBytes)); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		userID       string
		inputTestSQL string
		want         *User
		hasErr       bool
	}{
		{
			name:         "1件取得",
			userID:       "0001",
			inputTestSQL: "",
			want:         &User{"0001", "gopher1"},
			hasErr:       false,
		},
		{
			name:         "0件取得",
			userID:       "9999",
			inputTestSQL: "",
			want:         nil,
			hasErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlBytes, err := ioutil.ReadFile(tt.inputTestSQL)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := db.ExecContext(context.TODO(), string(sqlBytes)); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if _, err := db.ExecContext(context.TODO(), `TRUNCATE users;`); err != nil {
					t.Fatal(err)
				}
			})

			got, err := FetchUser(context.TODO(), tt.userID)
			if (err != nil) != tt.hasErr {
				t.Fatalf("FetchUser() error = %v, hasErr = %v", err, tt.hasErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchUser() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
