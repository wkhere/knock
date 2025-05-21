go:
	go build
	go test .

install: go
	go install

.PHONY: go install
