package configs

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"message-automation/src/models/base"
	"os"
	"strconv"
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
	messagePerExecution := 2 // Default
	if messagePerExecutionCustom, err := strconv.Atoi(os.Getenv("MESSAGE_PER_EXECUTION")); err == nil {
		messagePerExecution = messagePerExecutionCustom
	}

	executionPeriod := 2 // Default
	if executionPeriodCustom, err := strconv.Atoi(os.Getenv("EXECUTION_PERIOD")); err == nil {
		executionPeriod = executionPeriodCustom
	}
	return &Config{
		WebhookBaseURL:        "https://webhook.site",
		WebhookToken:          "ad14ab4c-d132-4ea5-a1b6-a8ded2932f02",
		MSSQLConnectionString: os.Getenv("DATABASE_URL"),
		RedisAddress:          os.Getenv("REDIS_ADDRESS"),
		RedisPassword:         "",
		RedisDB:               0,
		RedisTimeout:          time.Hour,
		MessagePerExecution:   messagePerExecution,
		ExecutionPeriod:       time.Minute * time.Duration(executionPeriod),
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
