all: genlogfmt

genlogfmt: *.go lib/*.go
	go build -o genlogfmt

test: *.go lib/*.go
	go test . ./lib
