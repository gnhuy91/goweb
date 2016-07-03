package main

import (
	"fmt"
	"os"
)

const dbDriver = "postgres"

func initDSN() string {
	// Prepare db DSN here
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "127.0.0.1:5432"
	}

	dsn := fmt.Sprintf("%s://%s:%s@%s/%s",
		"postgres",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		postgresHost,
		os.Getenv("POSTGRES_DB")) + "?sslmode=disable"

	return dsn
}
