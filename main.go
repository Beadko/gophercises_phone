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

	phoneNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, ph := range phoneNumbers {

		id, err := insertPhone(db, ph)
		if err != nil {
			fmt.Printf("Could not add a phone number: %v\n", err)
		} else {
			fmt.Println("Inserted phone with ID:", id)
		}
	}

	id := 2

	num, err := getPhone(db, id)
	if err != nil {
		fmt.Printf("Could not get the phone number for the id= %d: %v\n", id, err)
	} else {
		fmt.Printf("Phone number with ID %d: %s\n", id, num)
	}

}

func getPhone(db *sql.DB, id int) (string, error) {
	var num string
	if err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&num); err != nil {
		return "", err
	}
	return num, nil
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
