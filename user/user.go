package user

import (
	"pt-report-backend/db"
	"fmt"
	"github.com/golang-jwt/jwt/v4" // ✅ Correct package


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
// JWTCustomClaims struct
type JWTCustomClaims struct {
	UID       string `json:"uid"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Position  string `json:"position"`
	IsAdmin   bool   `json:"isAdmin"`
	jwt.StandardClaims
}

const userColumns = "uid, email, password, firstname, lastname"
// const (
// 	jwtKey            = "BdKPSNo7zxXR3P1h85klTMFWiaKP5KzbHO9A9bKcBAZ3xvknKAbYPmsrtaffFtJu"
// 	resetPwdJWTSecret = "TfleG7uvyzqXUOYu00Rzsnk87X49rvLyabTRacG8jtD58LeBc8e7y5gEjp7k48ku"
// )
const jwtKey = "BdKPSNo7zxXR3P1h85klTMFWiaKP5KzbHO9A9bKcBAZ3xvknKAbYPmsrtaffFtJu"


func GetJWTKey() []byte {
	return []byte(jwtKey)
}

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