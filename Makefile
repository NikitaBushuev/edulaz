GO := go

.PHONY: all

all: edulaz-d edulaz-cli

edulaz-d: cmd/d.go
	$(GO) build -o $@ $^

edulaz-cli: cmd/cli.go
	$(GO) build -o $@ $^
