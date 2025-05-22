package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "gophercises_phone"
)

var (
	errNotOpen = errors.New("Could not open the database:")
)

func main() {
	adminConnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable", user, password, host, port)

	db, err := sql.Open("pgx", adminConnStr)
	if err != nil {
		fmt.Printf("%v %v\n", errNotOpen, err)
		os.Exit(1)
	}

	if err = resetDB(db, dbname); err != nil {
		fmt.Printf("Could not create a database: %v\n", err)
		os.Exit(1)
	}
	db.Close()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)

	db, err = sql.Open("pgx", connStr)
	if err != nil {
		fmt.Printf("%v %v\n", errNotOpen, err)
		os.Exit(1)
	}
	defer db.Close()

	if err = createPhoneTable(db); err != nil {
		fmt.Printf("Could not create the table: %v\n", err)
		os.Exit(1)
	}

	id, err := insertPhone(db, "1234567890")
	if err != nil {
		fmt.Printf("Could not add a phone number: %v\n", err)
	}
	fmt.Println("Inserted phone with ID:", id)
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	stmt := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	if err := db.QueryRow(stmt, phone).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func createPhoneTable(db *sql.DB) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL PRIMARY KEY,
			value VARCHAR(255)
		)`
	_, err := db.Exec(stmt)
	return err
}

func resetDB(db *sql.DB, name string) error {
	if _, err := db.Exec(`DROP DATABASE IF EXISTS ` + quoteIdentifier(name)); err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec(`CREATE DATABASE ` + quoteIdentifier(name))
	return err
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}
