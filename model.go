package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Surname    string `json:"surname"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Postalcode string `json:"postalcode"`
	Confirmed  int    `json:"confirmed"`
	Verifycode string `json:"verifycode"`
	Regdate    string `json:"regdate"`
}

func (u *user) authUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE email='%s'", u.Email)
	pass := u.Password

	err := db.QueryRow(statement).Scan(&u.ID, &u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate)
	if err != nil {
		return err
	}
	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)); err != nil {
		return err
	}
	fmt.Println(err)
	return nil
}

func (u *user) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate)
}
func (u *user) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET name='%s', age=%d, surname='%s', email='%s', address='%s', city='%s', country='%s', postalcode='%s', confirmed=%d, verifycode='%s' WHERE id=%d", u.Name, u.Age, u.Surname, u.Email, u.Address, u.City, u.Country, u.Postalcode, u.Confirmed, u.Verifycode, u.ID)
	_, err := db.Exec(statement)
	return err
}
func (u *user) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}
func (u *user) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name FROM users WHERE email='%s'", u.Email)
	err := db.QueryRow(statement).Scan(&u.Name)
	switch err {
	case sql.ErrNoRows:
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			// TODO: Properly handle error
			log.Fatal(err)
		}
		statement = fmt.Sprintf("INSERT INTO users(name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate) VALUES('%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d,'%s','%s')", u.Name, u.Age, u.Surname, hash, u.Email, u.Address, u.City, u.Country, u.Postalcode, u.Confirmed, u.Verifycode, time.Now())
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)
		if err != nil {
			return err
		}
		statement = fmt.Sprintf("SELECT password, regdate FROM users WHERE id=%d", u.ID)
		err = db.QueryRow(statement).Scan(&u.Password, &u.Regdate)
		if err != nil {
			return err
		}

	default:

		err = errors.New("User Already exists")
		//log.Println("hi", err)
		return err

	}

	return nil
}
func getUsers(db *sql.DB, start, count int) ([]user, error) {
	statement := fmt.Sprintf("SELECT id, name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []user{}
	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}