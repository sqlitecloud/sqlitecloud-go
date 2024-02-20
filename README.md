
# SQLite Cloud Client SDK for Go v1.0.2

The SQLite Cloud Client SDK for Go (sqlitecloud/go-sdk) is the Go Programming Language application programmer's interface to [SQLite Cloud](https://sqlitecloud.io/). It is a set of library functions that allow client programs to pass queries and SQL commands to the SQLite Cloud backend server and to receive the results of these queries. In addition to the standard SQLite statements, several other [commands](https://docs.sqlitecloud.io/docs/commands) are supported.

## Getting Started

### Use the SQLite Cloud Client SDK in your Go code

1. Import the package in your Go source code:

    ```go
    import sqlitecloud "github.com/sqlitecloud/go-sdk"
    ```

2. Download the package, and run the [`go mod tidy` command](https://go.dev/ref/mod#go-mod-tidy) to synchronize your module's dependencies:

    ```
    $ go mod tidy 
    go: downloading github.com/sqlitecloud/go-sdk v1.0.0
    ```

3. Connect to a SQLite Cloud database with a valid [connection string](#get-a-connection-string):

    ```go
    db, err := sqlitecloud.Connect("sqlitecloud://user:pass@host.sqlite.cloud:port/dbname")
    ```

4. Execute queries using a [method](#api-documentation) defined on the `SQCloud` struct, for example `Select`:

    ```go
    result, _ := db.Select("SELECT * FROM table1;")
    ```

The following example shows how to print the content of the table `table1`:

```go
package main

import (
    "fmt"
    "strings"

    sqlitecloud "github.com/sqlitecloud/go-sdk"
)

const connectionString = "sqlitecloud://admin:password@host.sqlite.cloud:8860/dbname.sqlite"

func main() {
    db, err := sqlitecloud.Connect(connectionString)
    if err != nil {
        fmt.Println("Connect error: ", err)
    }

    tables, _ := db.ListTables()
    fmt.Printf("Tables:\n\t%s\n", strings.Join(tables, "\n\t"))

    fmt.Printf("Table1:\n")
    result, _ := db.Select("SELECT * FROM t1;")
    for r := uint64(0); r < result.GetNumberOfRows(); r++ {
        id, _ := result.GetInt64Value(r, 0)
        value, _ := result.GetStringValue(r, 1)
        fmt.Printf("\t%d: %s\n", id, value)
    }
}
```

## Get a connection string

You can connect to any cloud database using a special connection string in the form:

`sqlitecloud://user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3`

To get a valid connection string, follow these instructions:

- Get a [SQLite Cloud](https://sqlitecloud.io/) account. See the [documentation](https://docs.sqlitecloud.io/docs/introduction/login) for details.
- Create a [SQLite Cloud project](https://docs.sqlitecloud.io/docs/introduction/projects)
- Create a [SQLite Cloud database](https://docs.sqlitecloud.io/docs/introduction/databases)
- Get the connection string by clicking on the node address in the [Dashboard Nodes](https://docs.sqlitecloud.io/docs/introduction/nodes) section. A valid connection string will be copied to your clipboard.
- Add the database name to your connection string.

## API Documentation

The complete documentation of the sqlitecloud/go-sdk library is available at: https://pkg.go.dev/github.com/sqlitecloud/go-sdk