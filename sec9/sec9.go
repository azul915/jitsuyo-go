package sec9

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

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
	DROP TABLE IF EXISTS users;
	`)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		user_id varchar(32) NOT NULL,
		user_name varchar(100) NOT NULL,
		created_at timestamp with time zone,
		CONSTRAINT pk_users PRIMARY KEY (user_id)
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `
	INSERT INTO users VALUES 
	('0001','Gopher',current_timestamp),
	('0002','Ferris',current_timestamp)
	ON CONFLICT
	ON CONSTRAINT pk_users
	DO NOTHING;
	`); err != nil {
		log.Fatal(err)
	}

	type User struct {
		UserID    string
		UserName  string
		CreatedAt time.Time
	}

	rows, err := db.QueryContext(ctx, `
	SELECT user_id, user_name, created_at
	FROM users
	ORDER BY user_id;
	`)
	if err != nil {
		log.Fatalf("query all users: %v", err)
	}
	defer rows.Close()

	// rows, _ = db.Query("SELECT COUNT(*) as count FROM users")
	// var count int
	// for rows.Next() {
	// 	err := rows.Scan(&count)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// fmt.Println(count)

	var users []User
	for rows.Next() {
		var (
			userID, userName string
			createdAt        time.Time
		)

		if err := rows.Scan(&userID, &userName, &createdAt); err != nil {
			log.Fatalf("scan the user: %v", err)
		}
		fmt.Printf("userID: %v, userName: %v, createdAt: %v\n", userID, userName, createdAt)
		users = append(users, User{
			UserID:    userID,
			UserName:  userName,
			CreatedAt: createdAt,
		})
		if err := rows.Close(); err != nil {
			log.Fatalf("rows close: %v", err)
		}
		if err := rows.Err(); err != nil {
			log.Fatalf("scan users: %v", err)
		}
	}
	fmt.Println(users)

	var (
		userName  string
		createdAt time.Time
		userID    = "0002"
	)
	row := db.QueryRowContext(ctx, `
	SELECT user_name, created_at
	FROM users
	WHERE user_id = $1;
	`, userID)
	if err := row.Scan(&userName, &createdAt); err != nil {
		log.Fatalf("query row(user_id=%s): %v\n", userID, err)
	}
	fmt.Println(User{
		UserID:    userID,
		UserName:  userName,
		CreatedAt: createdAt,
	})

	// _, err = db.ExecContext(ctx, `
	// DROP TABLE IF EXISTS users;
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
