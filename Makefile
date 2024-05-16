deps:
	go mod download
	go mod verify
	go mod tidy

run: start

update_endpoint:
	curl -H "Token: test321" -X POST -f name=endpoint -f value=21 http://localhost:3000/endpoint

start:
	go run main.go ./contrib/dashboard.yaml
