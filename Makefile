all: logptn

logptn: *.go lib/*.go lib/*/*.go
	go build

test: *.go lib/*.go
	go test ./lib
