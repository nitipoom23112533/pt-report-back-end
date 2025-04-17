package main

import (
	"log"
	"pt-report-backend/api"
	"pt-report-backend/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := db.InitDB(db.Config{
		Username:     "pt-report",
		Password:     "6EuFmuTj3d83hgQ/",
		Server:       "35.186.157.84",
		DatabaseName: "pt_record",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())

	apiClient := api.NewAPI()
	apiClient.Group(e.Group(""))

	// เริ่มเซิร์ฟเวอร์
	e.Logger.Fatal(e.Start(":8080"))
}
