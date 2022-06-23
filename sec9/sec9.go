package sec9

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
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
	('0002','Ferris',current_timestamp),
	('0003','Duke',null)
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

	var (
		createdAtOrNull sql.NullTime
		dukeID          = "0003"
	)
	duke := db.QueryRowContext(ctx, `
	SELECT user_name, created_at
	FROM users
	WHERE user_id = $1;
	`, "0003")
	if err := duke.Scan(&userName, &createdAtOrNull); err != nil {
		log.Fatalf("query row(user_id=%s): %v\n", dukeID, err)
	}
	if !createdAtOrNull.Valid {
		createdAt = time.Time{}
	}
	fmt.Println(User{
		UserID:    dukeID,
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

type Service struct {
	db *sql.DB
}

func (s *Service) UpdateProduct(ctx context.Context, productID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `
		UPDATE products
		SET price = 200
		WHERE product_id = $1;
	`, productID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

type txAdmin struct {
	*sql.DB
}

type ProductService struct {
	tx txAdmin
}

func (t *txAdmin) Transaction(ctx context.Context, f func(ctx context.Context) (err error)) error {
	tx, err := t.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := f(ctx); err != nil {
		return fmt.Errorf("transaction query failed: %w", err)
	}
	return tx.Commit()
}

func (p *ProductService) UpdateProduct(ctx context.Context, productID string) error {
	updateFunc := func(ctx context.Context) error {
		if _, err := p.tx.ExecContext(ctx, `
		UPDATE products
		SET price = 200
		WHERE product_id = $1;
	`, productID); err != nil {
			return err
		}
		return nil
	}
	return p.tx.Transaction(ctx, updateFunc)
}

var _ pgx.Logger = (*logger)(nil)

type logger struct{}

func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if msg == "Query" {
		log.Printf("SQL:\n%v\nARGS:\n%v\n", data["sql"], data["args"])
	}
}

func Logging() {
	ctx := context.Background()
	config, err := pgx.ParseConfig("user=user password=password host=localhost port=5432 dbname=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	config.Logger = &logger{}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Exec(ctx, `
	DROP TABLE IF EXISTS users;`); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		user_id varchar(32) NOT NULL,
		user_name varchar(100) NOT NULL,
		created_at timestamp with time zone,
		CONSTRAINT pk_users PRIMARY KEY (user_id)
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Exec(ctx, `
	INSERT INTO users VALUES 
	('0001','Gopher',current_timestamp),
	('0002','Ferris',current_timestamp),
	('0003','Duke',current_timestamp)
	ON CONFLICT
	ON CONSTRAINT pk_users
	DO NOTHING;
	`); err != nil {
		log.Fatal(err)
	}

	sql := `SELECT schemaname, tablename FROM pg_tables WHERE schemaname = $1;`
	args := "information_schema"
	rows, err := conn.Query(ctx, sql, args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type PgTable struct {
		SchemaName string
		TableName  string
	}
	var pgtables []PgTable
	for rows.Next() {
		var s, t string
		if err := rows.Scan(&s, &t); err != nil {
			log.Fatal(err)
		}
		pgtables = append(pgtables, PgTable{SchemaName: s, TableName: t})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Print(pgtables)
}

func PreparedStatement() {

	// db, err := sql.Open("pgx", "host=localhost port=5432 user=user dbname=db password=password sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// ctx := context.Background()
	// if err := db.PingContext(ctx); err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = db.ExecContext(ctx, `
	// DROP TABLE IF EXISTS users;
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if _, err := db.ExecContext(ctx, `
	// CREATE TABLE IF NOT EXISTS users (
	// 	user_id varchar(32) NOT NULL,
	// 	user_name varchar(100) NOT NULL,
	// 	created_at timestamp with time zone,
	// 	CONSTRAINT pk_users PRIMARY KEY (user_id)
	// )`); err != nil {
	// 	log.Fatal(err)
	// }

	// type User struct {
	// 	UserID   string
	// 	UserName string
	// }
	// users := []User{
	// 	{"0001", "Gopher"},
	// 	{"0002", "Ferris"},
	// 	{"0003", "Duke"},
	// }

	// tx, err := db.BeginTx(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stmt, err := tx.PrepareContext(ctx, `
	// INSERT INTO users(
	// 	user_id,
	// 	user_name,
	// 	created_at
	// ) VALUES (
	// 	$1, $2, current_timestamp
	// );`)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()

	// for _, u := range users {
	// 	if _, err := stmt.ExecContext(ctx, u.UserID, u.UserName); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// if err := tx.Commit(); err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = db.ExecContext(ctx, `
	// TRUNCATE TABLE users;
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// valueStrings := make([]string, 0, len(users))
	// valueArgs := make([]interface{}, 0, len(users)*2)
	// number := 1
	// for _, u := range users {
	// 	valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", number, number+1))
	// 	valueArgs = append(valueArgs, u.UserID)
	// 	valueArgs = append(valueArgs, u.UserName)
	// 	number += 2
	// }
	// fmt.Printf("valueStrings: %v\n", valueStrings)
	// query := fmt.Sprintf(`
	// INSERT INTO users (
	// 	user_id, user_name
	// ) VALUES %s;
	// `, strings.Join(valueStrings, ","))
	// fmt.Printf("query: %v\n", query)
	// if _, err := db.ExecContext(ctx, query, valueArgs...); err != nil {
	// 	log.Fatal(err)
	// }

	// db.Close()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://user:password@localhost:5432/db")
	if err != nil {
		log.Fatal(err)
	}
	txn, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec(ctx, `
	DROP TABLE IF EXISTS products;
	`)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := txn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS products (
		product_no int2 NOT NULL,
		name varchar(32) NOT NULL,
		price int2 NOT NULL,
		CONSTRAINT pk_products PRIMARY KEY (product_no)
	);`); err != nil {
		log.Fatal(err)
	}

	rows := [][]interface{}{
		{1, "おにぎり", 120},
		{2, "パン", 300},
		{3, "お茶", 100},
	}
	_, err = txn.CopyFrom(ctx, pgx.Identifier{"products"}, []string{"product_no", "name", "price"}, pgx.CopyFromRows(rows))
	if err != nil {
		log.Fatal(err)
	}

}

type User struct {
	UserID   string
	UserName string
}

func FetchUser(ctx context.Context, userID string) (*User, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=user dbname=db password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row := db.QueryRowContext(ctx, `
	SELECT user_id, user_name
	FROM users
	WHERE user_id = $1;`, userID)
	user, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return user, nil
}

func scanUser(row *sql.Row) (*User, error) {
	var u User
	err := row.Scan(&u.UserID, &u.UserName)
	if err != nil {
		return nil, fmt.Errorf("row scan: %w", err)
	}
	return &u, nil
}
