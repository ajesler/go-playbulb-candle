default: build

build: clean
	go build -o pkg/candle-cli candle-cli/main.go

test:
	go test

clean:
	rm -rf build