run:
	@go run ./cmd/dummyapp/api -log-level DEBUG

test:
	@go test ./...

clean:
	@go mod tidy

build:
	@echo 'Building...'
	go build -o=./bin/dummyapp ./cmd/dummyapp/api
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/app ./cmd/dummyapp/api