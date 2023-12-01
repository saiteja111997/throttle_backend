clean:
	go clean
	rm -rf ./bin

build: clean
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

deploy_prod: clean build
	serverless deploy --stage prod --aws-profile saiteja

start: build
	sls offline --useDocker start --host 0.0.0.0 --stage local