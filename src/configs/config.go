package configs

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"message-automation/src/models/base"
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
	MessagePerExecution   int
	ExecutionPeriod       time.Duration
}

func SetConfig() *Config {
	return &Config{
		WebhookBaseURL:        "https://webhook.site",
		WebhookToken:          "",
		MSSQLConnectionString: "",
		RedisAddress:          "",
		RedisPassword:         "",
		RedisDB:               0,
		RedisTimeout:          time.Hour,
		MessagePerExecution:   2,
		ExecutionPeriod:       time.Minute * 2,
	}
}

func InitDB(dbConnectionString string) *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlserver", dbConnectionString)
	if err != nil {
		base.Log(fmt.Sprintf("Error while creating the connection pool %v", err.Error()))
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		base.Log(fmt.Sprintf("Error while pinging the database %v", err.Error()))
		os.Exit(1)
	}

	base.Log(fmt.Sprintf("MSSQL connection established"))
	return db
}
