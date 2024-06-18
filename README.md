
# Driver for SQLite Cloud

<p align="center">
  <img src="https://sqlitecloud.io/social/logo.png" height="300" alt="SQLite Cloud logo">
</p>

- [Driver for SQLite Cloud](#driver-for-sqlite-cloud)
- [Example](#example)
- [Get a Connection String](#get-a-connection-string)
- [Setting up the IDE](#setting-up-the-ide)

[![Test and QA](https://github.com/sqlitecloud/sqlitecloud-go/actions/workflows/testing.yaml/badge.svg?branch=main)](https://github.com/sqlitecloud/sqlitecloud-go/actions/workflows/testing.yaml)
[![codecov](https://codecov.io/gh/sqlitecloud/sqlitecloud-go/graph/badge.svg?token=5MAG3G4X01)](https://codecov.io/gh/sqlitecloud/sqlitecloud-go)
[![GitHub Tag](https://img.shields.io/github/v/tag/sqlitecloud/sqlitecloud-go?label=version&link=https%3A%2F%2Fpkg.go.dev%2Fgithub.com%2Fsqlitecloud%2Fsqlitecloud-go)](https://pkg.go.dev/github.com/sqlitecloud/sqlitecloud-go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sqlitecloud/sqlitecloud-go?link=https%3A%2F%2Fpkg.go.dev%2Fgithub.com%2Fsqlitecloud%2Fsqlitecloud-go)](https://pkg.go.dev/github.com/sqlitecloud/sqlitecloud-go)

---

[SQLite Cloud](https://sqlitecloud.io) for Go is a powerful package that allows you to interact with the SQLite Cloud database seamlessly. It provides methods for various database operations. This package is designed to simplify database operations in Go applications, making it easier than ever to work with SQLite Cloud. In addition to the standard SQLite statements, several other [commands](https://docs.sqlitecloud.io/docs/commands) are supported.

- Documentation: [https://pkg.go.dev/github.com/sqlitecloud/sqlitecloud-go#section-documentation](https://pkg.go.dev/github.com/sqlitecloud/sqlitecloud-go#section-documentation)
- Source: [https://github.com/sqlitecloud/sqlitecloud-go](https://github.com/sqlitecloud/sqlitecloud-go)
- Site: [https://sqlitecloud.io](https://sqlitecloud.io/developers)

## Example

### Use SQLite Cloud in your Go code

1. Import the package in your Go source code:

    ```go
    import sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
    ```

2. Download the package, and run the [`go mod tidy` command](https://go.dev/ref/mod#go-mod-tidy) to synchronize your module's dependencies:

    ```bash
    $ go mod tidy
    go: downloading github.com/sqlitecloud/sqlitecloud-go v1.0.0
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

    sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
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

## Get a Connection String

You can connect to any cloud database using a special connection string in the form:

`sqlitecloud://user:pass@host.com:port/dbname?timeout=10&key2=value2&key3=value3`

To get a valid connection string, follow these instructions:

- Get a [SQLite Cloud](https://sqlitecloud.io/) account. See the [documentation](https://docs.sqlitecloud.io/docs/introduction/login) for details.
- Create a [SQLite Cloud project](https://docs.sqlitecloud.io/docs/introduction/projects)
- Create a [SQLite Cloud database](https://docs.sqlitecloud.io/docs/introduction/databases)
- Get the connection string by clicking on the node address in the [Dashboard Nodes](https://docs.sqlitecloud.io/docs/introduction/nodes) section. A valid connection string will be copied to your clipboard.
- Add the database name to your connection string.



## Setting up the IDE

[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)

To start working on this project, follow these steps:

1. Open the project folder in Visual Studio Code (VSCode) using the remote container feature.
2. In the terminal, run the command `make setup-ide` to install the necessary development tools.
3. To ensure code quality, we have integrated [pre-commit](https://github.com/pre-commit/pre-commit) into the workflow. Before committing your changes to Git, pre-commit will run several tasks defined in the `.pre-commit-config.yaml` file.

By following these steps, you will have a fully set up development environment and be ready to contribute to the project.
