go:
	go build
	go test .

install: go
	go install -ldflags=-s

.PHONY: go install
