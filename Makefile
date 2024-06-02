export majorVersion=1
export minorVersion=0

export gittip=$(shell git log --format='%h' -n 1)
export patchVersion=$(shell git log --format='%h' | wc -l)
export ver=$(majorVersion).$(minorVersion).$(patchVersion).$(gittip)

include make/*.mk

deps:
	go mod download
	go mod verify
	go mod tidy

run: start

build:
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(ver)" -o build/dashboard main.go

start:
	go run main.go ./contrib/dashboard.yaml

binary: build
	./build/dashboard contrib/dashboard.yaml

.PHONY: build
