package auth

import (
	"pt-report-backend/user"
	"time"

	"github.com/golang-jwt/jwt"
)

const jwtKey = "BdKPSNo7zxXR3P1h85klTMFWiaKP5KzbHO9A9bKcBAZ3xvknKAbYPmsrtaffFtJu"

type JWTCustomClaims struct {
	UID       string `json:"uid"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	IsAdmin   bool   `json:"isAdmin"`
	jwt.StandardClaims
}

type Auth struct {
}

func GetJwtKey() []byte {
	return []byte(jwtKey)
}

func ParseJWTCustomClaims(x interface{}) *JWTCustomClaims {
	return x.(*jwt.Token).Claims.(*JWTCustomClaims)
}

func (s *Service) generateToken(x *user.User) (string, error) {
	claims := &JWTCustomClaims{
		UID:       x.UID,
		Email:     x.Email,
		Firstname: x.Firstname,
		Lastname:  x.Lastname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add((time.Hour * 24) * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
