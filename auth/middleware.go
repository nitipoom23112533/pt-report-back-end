package auth

import (
	"pt-report-backend/user"
    "github.com/labstack/echo/v4"
    echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/golang-jwt/jwt/v5"

)

func Auth() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey: []byte(user.GetJWTKey()),

		// ✅ ให้ข้าม OPTIONS method ไปเลย เพื่อไม่ให้โดน JWT block
		Skipper: func(c echo.Context) bool {
			return c.Request().Method == "OPTIONS"
		},
	}

	return echojwt.WithConfig(config)
}