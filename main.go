package main

import (
	"os"

	"github.com/fchoquet/bookmarks/app"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	users, err := app.ParseUsers(os.Getenv("BASIC_AUTH_USERS"))
	if err != nil {
		panic(err)
	}

	// no env vars should be accessed outside of the main function
	app.Start(app.Configuration{
		LogLevel:       os.Getenv("LOG_LEVEL"),
		BasicAuthUsers: users,
		DBConfig: app.DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Database: os.Getenv("DB_NAME"),
		},
		CSRFSecret:            []byte("sMKZudrjHxnN6fHLsJKFUMBeC7rnZ2Kd"),
		DisableCSRFProtection: os.Getenv("ENV") == "DEV",
	})
}
