build:
	CGO_ENABLED=0 go build -o sentinel ./cmd/sentinel/

run: build
	./sentinel

test:
	go test ./...

clean:
	rm -f sentinel

.PHONY: build run test clean
