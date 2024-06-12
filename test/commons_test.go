package test

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("./../.env"); err != nil {
		log.Print("No .env file found")
	}

	testConnectionString, _ := os.LookupEnv("SQLITE_CONNECTION_STRING")
	flag.StringVar(&testConnectionString, "server", testConnectionString, "Connection String")
}

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
