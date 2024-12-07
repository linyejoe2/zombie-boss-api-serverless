.PHONY: build clean deploy

build:
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bootstrap lambdas/main.go && "C:\Program Files\7-Zip\7z.exe" a bootstrap.zip bootstrap

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose


