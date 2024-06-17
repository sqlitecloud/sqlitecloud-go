module github.com/sqlitecloud/sqlitecloud-go/test

go 1.22

toolchain go1.22.3

require (
	github.com/joho/godotenv v1.5.1
	github.com/sqlitecloud/sqlitecloud-go v0.0.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xo/dburl v0.23.2 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sqlitecloud/sqlitecloud-go v0.0.0 => ../
