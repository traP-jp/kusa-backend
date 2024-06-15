package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

var (
	db  *sqlx.DB
	bot *traqwsbot.Bot
)

func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{
		User:                 os.Getenv("NS_MARIADB_USER"),
		Passwd:               os.Getenv("NS_MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("NS_MARIADB_HOSTNAME") + ":" + os.Getenv("NS_MARIADB_PORT"),
		DBName:               os.Getenv("NS_MARIADB_DATABASE"),
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	_db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("conntected")
	db = _db

	// bot, err = traqwsbot.NewBot(&traqwsbot.Options{
	// 	AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	e := echo.New()
	e.GET("/ping", pingHandler)
	e.GET("/tasks", tasksHandler)
	e.GET("/stamp/:id", getStampHandler)
	e.GET("/me", getMeHandler)
	e.GET("/rankings", getRankingsHandler)
	e.POST("/rankings", postRankingsHandler)

	e.Start(":8080")
}
