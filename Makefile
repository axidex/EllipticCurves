tidy:
	go fmt ./...
	go mod tidy

run:
	docker compose up -d --build

run-logs:
	docker compose up --build

run-logs-no-build:
	docker compose up
stop:
	docker compose down

re:
	docker compose down
	docker compose up -d --build

swag:
	swag init -g  ./cmd/main/main.go

update:
	go get github.com/axidex/Unknown

run-gui:
	go run cmd/gui/main.go

build:
	go env -w GOOS=windows
	go env -w GOARCH=amd64
	go build -ldflags "-H=windowsgui" -o win64.exe cmd/gui/main.go
