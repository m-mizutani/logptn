all: logptn

logptn: *.go lib/*.go
	go build -o logptn

test: *.go lib/*.go
	go test ./lib
