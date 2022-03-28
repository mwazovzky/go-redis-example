package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type User struct {
	Name  string
	Email string
}

func main() {
	db := connectDB()
	defer db.Close()

	users := getUsers(db)
	fmt.Println(users)

}

func connectDB() *sql.DB {
	dbhost := "postgres"
	dbport := 5432
	dbuser := "user"
	dbpassword := "secret"
	dbname := "testdb"

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected")

	return db
}

func getUsers(db *sql.DB) []User {
	users := []User{}
	var user User

	rows, err := db.Query("SELECT name, email FROM users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&user.Name, &user.Email)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users
}
