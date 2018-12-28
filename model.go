package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        string `json:"age"`
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
	Keeploged  int    `json:"keeploged"`
}

func (u *user) updateUserPass(db *sql.DB, id int, code string) error {
	log.Println(u.Password)
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	statement := fmt.Sprintf("SELECT id, name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE id = (select usersid from resetpass where id=%d and confirmed = 0)", id)
	err = db.QueryRow(statement).Scan(&u.ID, &u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	statement = fmt.Sprintf("UPDATE resetpass SET confirmed=1, newpass='%s' WHERE code='%s' and id=%d", hash, code, id)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	statement = fmt.Sprintf("UPDATE users SET password='%s' WHERE id=%d", hash, u.ID)
	log.Println("Statment:", statement)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	/* statement = fmt.Sprintf("SELECT name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE id =%d", u.ID)
	err = db.QueryRow(statement).Scan(&u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate) */
	return err
}
func (u *user) authUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE email='%s'", u.Email)
	pass := u.Password

	err := db.QueryRow(statement).Scan(&u.ID, &u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate)
	if err != nil {
		return err
	}
	if u.Confirmed != 1 {
		err = errors.New("User Not Confimed Yet Please Check Your Email Address")
		return err
	}
	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)); err != nil {
		return err
	}
	//fmt.Println(err)
	statement = fmt.Sprintf("INSERT INTO loginhist(userid, flag, logindate, keeploged) VALUES(%d, %d, '%s', %d)", u.ID, 1, time.Now(), u.Keeploged)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}
func (u *user) resetUser(db *sql.DB) error {

	statement := fmt.Sprintf("SELECT id, name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE email='%s'", u.Email)
	err := db.QueryRow(statement).Scan(&u.ID, &u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate)
	switch err {
	case sql.ErrNoRows:
		err = errors.New("User does not exist")
		return err
	default:

		code := randSeq(45)
		statement = fmt.Sprintf("INSERT INTO resetpass (usersid, code, confirmed, oldpass, resetdate) VALUES (%d, '%s', %d, '%s', '%s')", u.ID, code, 0, u.Password, time.Now())
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
		id := 0
		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&id)
		if err != nil {
			return err
		}

		mail.Send(u.Email, "Automated email from syncyours To Reset Password", "<strong>Reset Password: </strong><a href='http://"+config.Owner.URL+"/resetpass/"+strconv.Itoa(id)+"/"+code+"'>RESET</a>")

	}

	return nil
}
func (u *user) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate)
}
func (u *user) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET name='%s', age='%s', surname='%s', email='%s', address='%s', city='%s', country='%s', postalcode='%s', confirmed=%d, verifycode='%s' WHERE id=%d", u.Name, u.Age, u.Surname, u.Email, u.Address, u.City, u.Country, u.Postalcode, u.Confirmed, u.Verifycode, u.ID)
	_, err := db.Exec(statement)
	return err
}
func (u *user) verifyUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET confirmed=1 WHERE verifycode='%s' and id=%d", u.Verifycode, u.ID)
	_, err := db.Exec(statement)
	statement = fmt.Sprintf("SELECT name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate FROM users WHERE Verifycode='%s' and id =%d", u.Verifycode, u.ID)
	err = db.QueryRow(statement).Scan(&u.Name, &u.Age, &u.Surname, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.Postalcode, &u.Confirmed, &u.Verifycode, &u.Regdate)
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
		date := "%m/%d/%Y"
		statement = fmt.Sprintf("INSERT INTO users(name, age, surname, password, email, address, city, country, postalcode, confirmed, verifycode, regdate) VALUES('%s', STR_TO_DATE('%s', '%s'), '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d,'%s','%s')", u.Name, u.Age, date, u.Surname, hash, u.Email, u.Address, u.City, u.Country, u.Postalcode, u.Confirmed, u.Verifycode, time.Now())
		//log.Println(statement)
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)
		if err != nil {
			return err
		}
		verhash := randSeq(45)
		//log.Println(verhash)
		statement = fmt.Sprintf("UPDATE users SET verifycode = '%s' WHERE id=%d", verhash, u.ID)
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
		statement = fmt.Sprintf("SELECT password, verifycode, regdate FROM users WHERE id=%d", u.ID)
		log.Println(statement)
		err = db.QueryRow(statement).Scan(&u.Password, &u.Verifycode, &u.Regdate)
		if err != nil {
			return err
		}
		log.Println("URL: http://" + config.Owner.URL + "/user/" + strconv.Itoa(u.ID) + "/" + u.Verifycode + "")
		mail.Send(u.Email, "Automated email from syncyours", "<strong>test: </strong><a href='http://"+config.Owner.URL+"/verifyuser/"+strconv.Itoa(u.ID)+"/"+u.Verifycode+"'>Verify Account</a>")
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
