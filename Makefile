install:
	go mod tidy
build:
	GOOS=linux GOARCH=amd64 go build -o main ./src/main.go
test:
	go test ./src/... -v
deploy-commands:
	go run ./src/commands/commands.go -action create
deploy-infrastructure:
	cd infrastructure && cdk deploy
deploy:
	./bin/deploy