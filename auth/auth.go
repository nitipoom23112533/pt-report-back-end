package auth

import (
	"pt-report-backend/user"
	"time"

	"github.com/golang-jwt/jwt/v5"

)

const jwtKey = "BdKPSNo7zxXR3P1h85klTMFWiaKP5KzbHO9A9bKcBAZ3xvknKAbYPmsrtaffFtJu"


type JWTCustomClaims struct {
	UID       string `json:"uid"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	IsAdmin   bool   `json:"isAdmin"`
	jwt.RegisteredClaims

}
func (c *JWTCustomClaims) Valid() error {
	return nil
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
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add((time.Hour * 24) * 30)), // แปลงค่าเวลาเป็น NumericDate
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
