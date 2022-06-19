package sec9

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=user dbname=db password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		user_id varchar(32) NOT NULL,
		user_name varchar(100) NOT NULL,
		created_at timestamp with time zone,
		CONSTRAINT pk_users PRIMARY KEY (user_id)
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// _, err = db.ExecContext(ctx, `
	// DROP TABLE IF EXISTS users;
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
