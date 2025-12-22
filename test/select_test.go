package test

import (
	"os"
	"testing"

	"github.com/sqlitecloud/sqlitecloud-go"
)

func TestSelectWithLeadingChars(t *testing.T) {
	// Server API test
	connectionString, _ := os.LookupEnv("SQLITE_CONNECTION_STRING")
	apikey, _ := os.LookupEnv("SQLITE_API_KEY")
	connectionString += "/" + os.Getenv("SQLITE_DB") + "?apikey=" + apikey

	config, err1 := sqlitecloud.ParseConnectionString(connectionString)
	if err1 != nil {
		t.Fatal(err1.Error())
	}

	db := sqlitecloud.New(*config)
	err := db.Connect()

	if err != nil {
		t.Fatalf(err.Error())
	}

	defer db.Close()

	// test select with leading spaces
	result, err := db.Select("     SELECT 1 AS value;")
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test select with leading tabs
	result, err = db.Select("\t\t\tSELECT 1 AS value;")
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test select with leading new lines
	result, err = db.Select("\n\n\nSELECT 1 AS value;")
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test with leading carriage returns
	result, err = db.Select("\r\r\rSELECT 1 AS value;")
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test select with mixed leading characters
	result, err = db.Select(`
		SELECT 1 AS value;
	`)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}
}
