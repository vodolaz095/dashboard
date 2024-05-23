deps:
	go mod download
	go mod verify
	go mod tidy

run: start

update_endpoint:
	curl -v -H "Host: localhost" \
		-H "Token: test321" \
	    -H "Content-Type: application/x-www-form-urlencoded" \
		-X POST \
        -d "name=endpoint&value=21" \
		http://localhost:3000/update

start:
	go run main.go ./contrib/dashboard.yaml
