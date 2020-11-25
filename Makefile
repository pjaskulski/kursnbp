build:
	go build -o kursnbp ./cmd/kursnbp

run:
	go run ./cmd/kursnbp gold

linux: 
	go build -ldflags="-s -w" -o builds/linux/kursnbp ./cmd/kursnbp

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o builds/windows/kursnbp.exe ./cmd/kursnbp

macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/macos/kursnbp ./cmd/kursnbp

freebsd:
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o builds/freebsd/kursnbp ./cmd/kursnbp

all: linux windows macos freebsd

test:
	go test -v ./pkg/nbpapi

testcheck:
	go test -v -run TestCheckArg ./pkg/nbpapi

cover:
	go test -cover ./pkg/nbpapi

release:
	sh scripts/release.sh

