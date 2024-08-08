export majorVersion=1
export minorVersion=0

export gittip=$(shell git log --format='%h' -n 1)
export patchVersion=$(shell git log --format='%h' | wc -l)
export ver=$(majorVersion).$(minorVersion).$(patchVersion).$(gittip)

include make/*.mk

tools:
	@which podman
	@podman version
	@which redis-cli
	@redis-cli --version
	@which go
	@go version

# https://go.dev/blog/govulncheck
# install it by go install golang.org/x/vuln/cmd/govulncheck@latest
vuln:
	which govulncheck
	govulncheck ./...

deps:
	go mod download
	go mod verify
	go mod tidy

test:
	go test -v ./...

run: start

build: deps
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(ver)" -o build/dashboard main.go

build/linux_amd64: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(ver)" -o build/dashboard_linux_amd64 main.go

build/linux_arm6: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-X main.Version=$(ver)" -o build/dashboard_linux_arm6 main.go

build/linux_arm7: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main.Version=$(ver)" -o build/dashboard_linux_arm7 main.go

build/windows: deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(ver)" -o build/dashboard.exe main.go

build/macos: deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(ver)" -o build/dashboard_darwin_amd64 main.go

build/all: build/linux_amd64 build/linux_arm6 build/linux_arm7 build/windows build/macos
	md5sum build/dashboard* > build/dashboard.md5
	ls -hl build/

start:
	go run main.go ./contrib/dashboard.yaml

binary: build
	./build/dashboard contrib/dashboard.yaml

tag:
	git tag "v$(majorVersion).$(minorVersion).$(patchVersion)"

.PHONY: build
