run: build
	@./bin/blueMQ

build:
	@go build -o bin/blueMQ
	