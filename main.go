package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := ConnectToDB()
	if err != nil {
		log.Fatalf("failed connecting to db %s", err.Error())
	}
	rand.Seed(time.Now().Unix())
	max := 50
	if err := PopulateUsers(max, db); err != nil {
		fmt.Printf("Failed populating users | %s", err.Error())
	}

	max = 70
	if err := PopulateCreditCards(max, db); err != nil {
		fmt.Printf("Failed populating credit cards | %s\n", err.Error())
	}
	fmt.Println(db)
}

func ConnectToDB() (*sql.DB, error) {
	return sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
}

func PopulateUsers(max int, db *sql.DB) error {
	firstnames, err := ReadFileToStringArray("lists/firstnames.txt")
	if err != nil {
		return fmt.Errorf("failed creating firstnames list | %s", err.Error())
	}
	lastnames, err := ReadFileToStringArray("lists/lastnames.txt")
	if err != nil {
		return fmt.Errorf("failed creating lastnames list | %s", err.Error())
	}
	cities, err := ReadFileToStringArray("lists/cities.txt")
	if err != nil {
		return fmt.Errorf("failed creating cities list | %s", err.Error())
	}
	passwords, err := ReadFileToStringArray("lists/passwords.txt")
	if err != nil {
		return fmt.Errorf("failed creating passwords list | %s", err.Error())
	}

	for i := 0; i < max; i++ {
		firstname := pick(firstnames)
		lastname := pick(lastnames)
		city := pick(cities)
		password := pick(passwords)

		if _, err := db.Query(
			fmt.Sprintf(
				"insert into wp_client_infos (Nom, Prenom, Adresse, Email, Password)"+
					"values (%s, %s, %s, %s)",
				firstname, lastname, city, password,
			),
		); err != nil {
			fmt.Printf("Failed inserting row | %s\n", err.Error())
		}

	}
	return nil
}

func ReadFileToStringArray(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(content), "\n"), nil
}

func pick(in []string) string {
	return in[rand.Intn(len(in))]
}

func PopulateCreditCards(max int, db *sql.DB) error {
	return nil
}
