package postgres

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	_ "github.com/lib/pq"
)

func DBInit(logger *zap.Logger) (*sqlx.DB) {
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	database := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal("cannot connect to database ", zap.Error(err))
	}

	return db
}