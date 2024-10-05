package configs

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"os"
	"time"
)

type Config struct {
	WebhookBaseURL        string
	WebhookToken          string
	MSSQLConnectionString string
	RedisAddress          string
	RedisPassword         string
	RedisDB               int
	RedisTimeout          time.Duration
}

func SetConfig() *Config {
	return &Config{
		WebhookBaseURL:        "https://webhook.site",
		WebhookToken:          "fc8ec7bc-5939-43aa-a1d7-a85bbe6f1a13",
		MSSQLConnectionString: os.Getenv("DATABASE_URL"),
		RedisAddress:          "localhost:6379",
		RedisPassword:         "",
		RedisDB:               0,
		RedisTimeout:          time.Hour,
	}
}

func InitDB(dbConnectionString string) *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlserver", dbConnectionString)
	if err != nil {
		log.Fatalf("Error occurred while creating connection pool: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error while pinging the database: %s", err.Error())
	}

	fmt.Printf("\n[%s] MSSQL connection established", time.Now().Format("2006-01-02.15.04.05"))
	return db
}
