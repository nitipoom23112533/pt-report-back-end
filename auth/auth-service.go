package auth

import (
	"pt-report-backend/user"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/golang-jwt/jwt"
	"time"
	"github.com/labstack/echo/v4"


)

type Service struct {
	userService *user.Service
}

func NewService(userService *user.Service) *Service {
	return &Service{userService: userService}
}

var (
	ErrLogin = errors.New("email or password invalid")
)

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrLogin
		}
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrLogin
	}
	return s.generateToken(user)
}


var (
	ErrInvalidToken = errors.New("invalid token")
)
// VerifyToken ตรวจสอบว่า Token ถูกต้องหรือไม่
func (s *Service) VerifyToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJwtKey(), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	// ตรวจสอบว่า Token ถูกต้องและยังไม่หมดอายุ
	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// ตรวจสอบว่า Token หมดอายุหรือไม่
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}


// // VerifyTokenHandler สำหรับตรวจสอบ Token
// func VerifyTokenHandler(c echo.Context) error {
// 	tokenString := c.QueryParam("token") // ดึง Token จาก Query Param

// 	if tokenString == "" {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
// 	}

// 	// แปลง token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return GetJwtKey(), nil
// 	})

// 	if err != nil || !token.Valid {
// 		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
// 	}

// 	// ดึง Claims
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
// 	}

// 	// ส่งข้อมูลผู้ใช้กลับไป
// 	return c.JSON(http.StatusOK, claims)
// }

func VerifyTokenHandler(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization") // ดึง Token จาก Header: Authorization

	// ตรวจสอบว่า Token มีหรือไม่
	if tokenString == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}

	// ถ้ามี Bearer token ก็ให้ตัดคำว่า "Bearer " ออก
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// แปลง token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return GetJwtKey(), nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	// ดึง Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}

	// ส่งข้อมูลผู้ใช้กลับไป
	return c.JSON(http.StatusOK, claims)
}

