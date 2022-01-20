.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/converter ./converter
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/shots-retreiver ./shots-retreiver

clean:
	rm -rf ./bin

deploy-converter: clean build
	sls deploy function -f converter
	sls deploy function -f shotsRetreiver
