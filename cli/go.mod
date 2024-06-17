module github.com/sqlitecloud/go-cli

go 1.22

toolchain go1.22.3

require (
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/peterh/liner v1.2.2
	golang.org/x/term v0.21.0
)

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-runewidth v0.0.3 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/xo/dburl v0.23.2 // indirect
	golang.org/x/sys v0.21.0 // indirect
)

require github.com/sqlitecloud/sqlitecloud-go v0.0.0

replace github.com/sqlitecloud/sqlitecloud-go v0.0.0 => ../
