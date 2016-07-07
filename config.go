package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gnhuy91/go-vcap-parser"
)

const dbDriver = "postgres"

var (
	uaaURI           = getUaaURI()
	uaaCheckTokenURI = uaaURI + "/check_token"
)

func configPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func configDSN() string {
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
			log.Println("Error reading VCAP_SERVICES env:", err)
			return dsn
		}

		pg, prs := vcap["postgres"]
		if !prs {
			log.Println(`Error reading "postgres" from VCAP_SERVICES`)
			return dsn
		}
		if len(pg) == 0 {
			log.Println(`Error reading "postgres" from VCAP_SERVICES: index out of range`)
			return dsn
		}
		dsn = pg[0].Credentials.DSN
	}

	return dsn
}

func getUaaURI() string {
	// Check if UAA_URI env is provided
	if uri := os.Getenv("UAA_URI"); uri != "" {
		return uri
	}

	// If no UAA_URI provided, try parsing it from VCAP_SERVICES
	vcapServices := os.Getenv("VCAP_SERVICES")
	if vcapServices == "" {
		return ""
	}
	vcap, err := vcapparser.ParseVcapServices(vcapServices)
	if err != nil {
		log.Println("Error reading VCAP_SERVICES env:", err)
		return ""
	}

	uaa, prs := vcap["predix-uaa"]
	if !prs {
		log.Println(`Error reading "predix-uaa" from VCAP_SERVICES`)
		return ""
	}
	if len(uaa) == 0 {
		log.Println(`Error reading "predix-uaa" from VCAP_SERVICES: index out of range`)
		return ""
	}

	return uaa[0].Credentials.URI
}
