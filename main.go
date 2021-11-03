package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"strconv"
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
	domains, err := ReadFileToStringArray("lists/domains.txt")
	if err != nil {
		return fmt.Errorf("failed creating domains list | %s", err.Error())
	}

	for i := 0; i < max; i++ {
		firstname := strings.ToLower(Pick(firstnames))
		lastname := Pick(lastnames)
		city := Pick(cities)
		password := Pick(passwords)
		email := fmt.Sprintf("%s@%s", generateAliasName([]string{firstname, lastname, city, password}), Pick(domains))

		if _, err := db.Query(
			fmt.Sprintf(
				"insert into wp_client_infos (Nom, Prenom, Adresse, Email, Password)"+
					"values ('%s', '%s', '%s', '%s', '%s')",
				lastname, firstname, city, email, password,
			),
		); err != nil {
			fmt.Printf("%s | %s | %s | %s | %s \n", firstname, lastname, city, email, password)
			fmt.Printf("Failed inserting row | %s\n", err.Error())
		}
	}
	return nil
}

func generateAliasName(infos []string) string {
	var aliases []string
	for _, info := range infos {

		aliases = append(aliases, strings.TrimSpace(DeleteSpecialChars(strings.ToLower(info[:RandRange(0, len(info))]))))

	}
	rand.Shuffle(len(aliases), func(i, j int) { aliases[i], aliases[j] = aliases[j], aliases[i] })
	return strings.Join(aliases, "")
}

func DeleteSpecialChars(target string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(target, "")
}

func ReadFileToStringArray(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(content), "\n"), nil
}

func Pick(in []string) string {
	return in[rand.Intn(len(in))]
}

func RandRange(min int, max int) int {
	if max == 0 {
		max = 1
	}
	return rand.Intn(max-min) + min
}

func PopulateCreditCards(max int, db *sql.DB) error {
	firstnames, err := ReadFileToStringArray("lists/firstnames.txt")
	if err != nil {
		return fmt.Errorf("failed creating firstnames list | %s", err.Error())
	}
	lastnames, err := ReadFileToStringArray("lists/lastnames.txt")
	if err != nil {
		return fmt.Errorf("failed creating lastnames list | %s", err.Error())
	}
	for i := 0; i < max; i++ {
		date := RandomDate()
		number := RandomCardNumber()
		cvv := RandomCVV()

		if _, err := db.Query(
			fmt.Sprintf(
				"insert into wp_credit_card (Number, Expiration, CVV, Name)"+
					"values ('%s', '%s', '%s', '%s')",
				number, date, cvv, fmt.Sprintf("%s %s ", strings.ToUpper(Pick(firstnames)), strings.ToUpper(Pick(lastnames))),
			),
		); err != nil {
			fmt.Printf("Failed inserting row | %s\n", err.Error())
		}
	}
	return nil
}

func RandomDate() string {
	month := RandRange(1, 12)
	year := RandRange(20, 30)

	return fmt.Sprintf("%s/%s", fmt.Sprintf("%02d", month), strconv.Itoa(year))
}

func RandomCardNumber() string {
	var numbers []string
	for i := 0; i < 4; i++ {
		nb := RandRange(0, 9999)
		numbers = append(numbers, fmt.Sprintf("%04d", nb))
	}
	return strings.Join(numbers, "-")
}

func RandomCVV() string {
	return fmt.Sprintf("%03d", RandRange(0, 999))
}
