run:
	@go run ./cmd/dummy-app/api -log-level DEBUG

test:
	@go test ./...

clean:
	@go mod tidy

build:
	@echo 'Building...'
	go build -o=./bin/dummy-app ./cmd/dummy-app/api
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/app ./cmd/dummy-app/api