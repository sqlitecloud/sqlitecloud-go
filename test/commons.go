package test

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
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

func setupDatabase(t *testing.T) (*sqlitecloud.SQCloud, func()) {
	connectionString, _ := os.LookupEnv("SQLITE_CONNECTION_STRING")
	apikey, _ := os.LookupEnv("SQLITE_API_KEY")

	db, err := sqlitecloud.Connect(connectionString + "?apikey=" + apikey)
	if err != nil {
		t.Fatal("Connection error: ", err.Error())
	}

	return db, func() {
		db.Close()
	}
}
