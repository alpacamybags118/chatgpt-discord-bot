install:
	go mod tidy
build:
	go build -o ./tmp/main src/main.go
test:
	go test ./src/... -v
deploy-infrastructure:
	cd infrastructure && cdk deploy
deploy:
	./bin/deploy