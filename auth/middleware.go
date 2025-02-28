package auth

import (
	"pt-report-backend/user"
    "github.com/labstack/echo/v4"
    echojwt "github.com/labstack/echo-jwt/v4"

)

func Auth() echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey: []byte(user.GetJWTKey()), // ใส่ JWT Secret Key
		TokenLookup: "header:Authorization",  // กำหนดตำแหน่ง JWT Token
		// AuthScheme:  "Bearer",               // Bearer Token
	}
	return echojwt.WithConfig(config)
}