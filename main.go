package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/lib/pq"
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
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("%v %v\n", errNotOpen, err)
		os.Exit(1)
	}

	if err = resetDB(db, dbname); err != nil {
		fmt.Printf("Could not create a database: %v\n", err)
		os.Exit(1)
	}
	db.Close()*/
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
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
	fmt.Println(id)
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
			id SERIAL,
			value VARCHAR(255)
		)`
	_, err := db.Exec(stmt)
	return err
}

func resetDB(db *sql.DB, name string) error {
	if _, err := db.Exec("DROP DATABASE IF EXISTS " + name); err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	valid := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !valid.MatchString(name) {
		return errors.New("invalid database name")
	}
	if _, err := db.Exec("CREATE DATABASE " + pq.QuoteIdentifier(name)); err != nil {
		return err
	}
	return nil
}

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

/* func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}*/
