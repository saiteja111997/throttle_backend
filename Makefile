clean:
	go clean
	rm -rf ./bin/fileUploadService

build fileUploadService: clean
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/fileUploadService/main fileUploadService/main.go

deploy_prod: clean build fileUploadService
	serverless deploy --stage prod --aws-profile saiteja

start: build
	sls offline --useDocker start --host 0.0.0.0 --stage local