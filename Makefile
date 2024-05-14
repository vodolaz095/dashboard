deps:
	go mod download
	go mod verify
	go mod tidy

run: start

start:
	go run main.go ./contrib/dashboard.yaml
