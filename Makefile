build:
	GOOS=linux go build -o ohlc cmd/ohlc/*.go

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
