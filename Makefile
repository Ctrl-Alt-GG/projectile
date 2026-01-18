export GOOS := linux
export GOARCH := amd64
export CGO_ENABLED := 0

.PHONY:
all: bin/server bin/agent

.PHONY:
bin/server: pkg/agentmsg/agentmsg.pb.go pkg/agentmsg/agentmsg_grpc.pb.go bin
	go build -v -o bin/server ./cmd/server

.PHONY:
bin/agent: pkg/agentmsg/agentmsg.pb.go pkg/agentmsg/agentmsg_grpc.pb.go bin
	go build -v -o bin/agent ./cmd/agent

bin:
	mkdir -v bin

pkg/agentmsg:
	mkdir -pv pkg/agentmsg

pkg/agentmsg/agentmsg.pb.go pkg/agentmsg/agentmsg_grpc.pb.go: agentmsg.proto pkg/agentmsg
	protoc -I. --go_out=./pkg/agentmsg --go_opt=paths=source_relative --go-grpc_out=./pkg/agentmsg --go-grpc_opt=paths=source_relative agentmsg.proto