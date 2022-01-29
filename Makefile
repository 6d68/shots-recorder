.PHONY: build clean deploy test

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/converter ./converter
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/shots-retreiver ./shots-retreiver
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/shots-url-signer ./shots-url-signer

clean:
	rm -rf ./bin

test:
	go test ./...

deploy-converter: clean build
	sls deploy function -f converter
	sls deploy function -f shotsRetreiver
	sls deploy function -f shotsUrlSigner
