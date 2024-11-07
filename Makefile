.PHONY: build clean deploy

build:
	# dep ensure -v
	# env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	# env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bootstrap lambda/main.go && "C:\Program Files\7-Zip\7z.exe" a bootstrap.zip bootstrap

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose


