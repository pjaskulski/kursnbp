compile:
	go build -o kursnbp ./cmd/kursnbp

run:
	go run ./cmd/kursnbp gold

linux: 
	go build -ldflags="-s -w" -o kursnbp ./cmd/kursnbp

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/windows/kursnbp.exe ./cmd/kursnbp

macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/macos/kursnbp ./cmd/kursnbp

freebsd:
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o build/freebsd/kursnbp ./cmd/kursnbp

all: linux windows macos freebsd

test:
	go test -v ./cmd/kursnbp

testcheck:
	go test -v -run TestCheckArg ./cmd/kursnbp	

cover:
	go test -cover ./cmd/kursnbp
