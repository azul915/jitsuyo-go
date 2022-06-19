package sec9

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=user dbname=db password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(db.Stats())
}
