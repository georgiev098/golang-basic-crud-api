package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/georgiev098/golang-basic-crud-api/pkg/utils"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDB(dbName string) (*sql.DB, error) {
	utils.LoadEnv()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	// Build DSN string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, name)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Pool settings
	maxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	lifetime, _ := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))

	if maxOpen > 0 {
		db.SetMaxOpenConns(maxOpen)
	}
	if maxIdle > 0 {
		db.SetMaxIdleConns(maxIdle)
	}
	if lifetime > 0 {
		db.SetConnMaxLifetime(lifetime)
	}

	// Actual connection test
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("✅ Successfully connected to MariaDB")
	return db, nil
}
