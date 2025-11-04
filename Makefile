export GOOS := linux
export GOARCH := amd64
export CGO_ENABLED := 0

bin/projectile: model/model.pb.go main.go bin
	go build -v -o bin/projectile .

bin:
	mkdir -v bin

model/model.pb.go:
	protoc -I. --go_out=./model --go_opt=paths=source_relative model.proto