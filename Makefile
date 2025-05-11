export GOOS := linux
export GOARCH := amd64
export CGO_ENABLED := 0

bin/projectile: main.go bin
	go build -v -o bin/projectile .

bin:
	mkdir -v bin