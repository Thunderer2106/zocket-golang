package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool" 
)

var TestDB *pgxpool.Pool 

func SetupTestDB() *pgxpool.Pool {
	connStr := "postgresql://postgres:postgres@localhost:5432/pgdb"

	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to test database: %v\n", err)
	}
	TestDB = pool
	return pool
}

func TearDownTestDB() {
	if TestDB != nil {
		TestDB.Close()
	}
}
