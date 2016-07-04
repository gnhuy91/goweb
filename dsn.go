package main

import (
	"fmt"
	"os"

	"github.com/gnhuy91/go-vcap-parser"
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

	// For deploying to Cloud Foundry
	vcapServices := os.Getenv("VCAP_SERVICES")
	if vcapServices != "" {
		vcap, err := vcapparser.ParseVcapServices(vcapServices)
		if err != nil {
			fmt.Println(err)
		}
		dsn = vcap["postgres"][0].Credentials.DSN
	}

	return dsn
}
