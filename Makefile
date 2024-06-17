GOPATH		= $(shell go env GOPATH)

# Test SDK
test:
	cd test; go test -v

# GO SDK
sdk:	*.go
	go build .

# CLI App
$(GOPATH)/bin/sqlc:	*.go cli/sqlc.go
	cd cli; go build -o $(GOPATH)/bin/sqlc

cli: $(GOPATH)/bin/sqlc

github:
	open https://github.com/sqlitecloud/sqlitecloud-go

diff:
	git difftool


# gosec
gosec:
ifeq ($(wildcard $(GOPATH)/bin/gosec),)
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
endif

checksec:	gosec *.go cli/*.go
	$(GOPATH)/bin/gosec -exclude=G304 .
	cd cli; $(GOPATH)/bin/gosec -exclude=G304,G302 .


# Documentation
godoc:
ifeq ($(wildcard $(GOPATH)/bin/godoc),)
	go install golang.org/x/tools/cmd/godoc
endif

doc:	godoc
ifeq ($(wildcard ./src),)
	ln -s . src
endif
	@echo "Hit CRTL-C to stop the documentation server..."
	@( sleep 1 && open http://localhost:6060/pkg/github.com/sqlitecloud/sqlitecloud-go/ ) &
	@$(GOPATH)/bin/godoc -http=:6060 -index -play

clean:
	rm -rf $(GOPATH)/bin/sqlc*

all: sdk cli test

.PHONY: test sdk
