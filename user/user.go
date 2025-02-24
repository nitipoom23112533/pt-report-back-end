package user

import (
	"pt-report-backend/db"
	"fmt"
)

type User struct {
	UID       string `db:"uid" json:"uid"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
}

type Role struct {
	ID   string `db:"id" json:"id"`
	UID  string `db:"uid" json:"uid"`
	Role string `db:"role" json:"role"`
}

const userColumns = "uid, email, password, firstname, lastname"

func getUsers() ([]User, error) {
	query := fmt.Sprintf(`SELECT %s
							FROM users`, userColumns)
	var xs []User
	err := db.DB.Select(&xs, query)
	return xs, err
}

func getUserByEmail(email string) (*User, error) {
	query := fmt.Sprintf(`SELECT %s
							FROM users
							WHERE email = ?
							LIMIT 1`, userColumns)
	var x User
	err := db.DB.Get(&x, query, email)
	return &x, err
}

func getRoleByUID(uid string) (*Role, error) {
	query := `select
							id, uid, role
					from role 
					where uid = ?`
	var x Role
	err := db.DB.Get(&x, query, uid)
	return &x, err
}