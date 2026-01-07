//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.0
//     //             ///   ///  ///    Date        : 2021/10/08
//    ///             ///   ///  ///    Author      : Andrea Donetti
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Test program for the
//   ////                ///  ///                     SQLite Cloud internal
//     ////     //////////   ///                      server commands.
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package test

import (
	"os"
	"testing"

	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
)

const testDbnameSelectArray = "test-gosdk-selectarray-db.sqlite"

func TestSelectArray(t *testing.T) {
	// Server API test
	connectionString, _ := os.LookupEnv("SQLITE_CONNECTION_STRING")
	apikey, _ := os.LookupEnv("SQLITE_API_KEY")
	connectionString += "?apikey=" + apikey

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

	// test select null string, it was causing a crash on the server
	if res, err := db.Select(""); err != nil {
		t.Fatal(err.Error())
	} else if !res.IsNULL() {
		t.Fatalf("Expected NULL, got %v", res.GetType())
	}

	// Checking CREATE DATABASE
	if err := db.ExecuteArray("CREATE DATABASE ? PAGESIZE ? IF NOT EXISTS", []interface{}{testDbnameSelectArray, 4096}); err != nil {
		t.Fatal(err.Error())
	}

	// Checking USE DATABASE
	if err := db.UseDatabase(testDbnameSelectArray); err != nil { // Database
		t.Fatal(err.Error())
	}

	// Creating Table
	if err := db.Execute("CREATE TABLE IF NOT EXISTS t1 (a INTEGER PRIMARY KEY, b)"); err != nil {
		t.Fatal(err.Error())
	}

	// Deleting Table content
	if err := db.Execute("DELETE FROM t1"); err != nil {
		t.Fatal(err.Error())
	}

	// Adding rows to Table with ExecuteArray
	if err := db.ExecuteArray("INSERT INTO t1 (b) VALUES (?), (?), (?), (?)", []interface{}{int(1), "text 2", 2.2, []byte("A")}); err != nil {
		t.Fatal(err.Error())
	}

	// Select Table
	if res, err := db.SelectArray("SELECT * FROM t1 WHERE a >= ?", []interface{}{1}); err != nil {
		t.Fatal(err.Error())
	} else if res.GetNumberOfRows() != 4 {
		t.Fatalf("Expected 4 rows, got %d", res.GetNumberOfRows())
	} else if res.GetNumberOfColumns() != 2 {
		t.Fatalf("Expected 2 columns, got %d", res.GetNumberOfColumns())
	} else if s, _ := res.GetStringValue(0, 1); s != "1" {
		t.Fatalf("Expected '1', got '%s'", s)
	} else if s, _ := res.GetStringValue(1, 1); s != "text 2" {
		t.Fatalf("Expected 'text 2', got '%s'", s)
	} else if s, _ := res.GetStringValue(2, 1); s != "2.2" {
		t.Fatalf("Expected '2.2', got '%s'", s)
	} else if s, _ := res.GetStringValue(3, 1); s != "A" {
		t.Fatalf("Expected 'A', got '%s'", s)
	}

	// Checking Unuse Database
	if err := db.UnuseDatabase(); err != nil {
		t.Fatal(err.Error())
	}

	// Checking REMOVE DATABASE
	if err := db.RemoveDatabase(testDbnameSelectArray, false); err != nil {
		t.Fatal(err.Error())
	}
}

func TestSelectArrayTableNameWithQuotes(t *testing.T) {
	// Server API test
	connectionString, _ := os.LookupEnv("SQLITE_CONNECTION_STRING")
	apikey, _ := os.LookupEnv("SQLITE_API_KEY")
	connectionString += "?apikey=" + apikey

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

	// Checking CREATE DATABASE
	if err := db.ExecuteArray("CREATE DATABASE ? IF NOT EXISTS", []interface{}{testDbnameSelectArray}); err != nil {
		t.Fatal(err.Error())
	}

	// Creating Table
	if err := db.ExecuteArray("SWITCH DATABASE ?; CREATE TABLE 'table '' name' (id INTEGER PRIMARY KEY, value)", []interface{}{testDbnameSelectArray}); err != nil {
		t.Fatal(err.Error())
	}

	// Checking REMOVE DATABASE
	if err := db.RemoveDatabase(testDbnameSelectArray, false); err != nil {
		t.Fatal(err.Error())
	}
}

func TestSelectArrayWithLeadingChars(t *testing.T) {
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
	result, err := db.SelectArray("     SELECT 1 AS value;", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test select with leading tabs
	result, err = db.SelectArray("\t\t\tSELECT 1 AS value;", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test select with leading new lines
	result, err = db.SelectArray("\n\n\nSELECT 1 AS value;", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test with leading carriage returns
	result, err = db.SelectArray("\r\r\rSELECT 1 AS value;", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}

	// test select with mixed leading characters
	result, err = db.SelectArray(`
		SELECT 1 AS value;
	`, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result.GetNumberOfRows() != 1 {
		t.Fatalf("Expected 1 row, got %d rows", result.GetNumberOfRows())
	}
}
