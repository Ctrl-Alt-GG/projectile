export GOOS := linux
export GOARCH := amd64
export CGO_ENABLED := 1

BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
LATEST_TAG := $(shell git describe --tags --abbrev=0)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
DIRTY := $(shell git diff-index --quiet HEAD -- || echo "x")

LDFLAGS := -X 'github.com/Ctrl-Alt-GG/projectile/pkg/utils.version=${LATEST_TAG}' -X 'github.com/Ctrl-Alt-GG/projectile/pkg/utils.commitHash=${COMMIT_HASH}' -X 'github.com/Ctrl-Alt-GG/projectile/pkg/utils.buildTimestamp=${BUILD_TIMESTAMP}' -X 'github.com/Ctrl-Alt-GG/projectile/pkg/utils.dirty=${DIRTY}'

.PHONY:
all: bin/server bin/agent

.PHONY:
clean:
	rm -rf bin/

bin/server: pkg/agentmsg/agentmsg.pb.go pkg/agentmsg/agentmsg_grpc.pb.go | bin
	go build -v -buildvcs=false -ldflags="${LDFLAGS}" -o bin/server ./cmd/server

bin/agent: pkg/agentmsg/agentmsg.pb.go pkg/agentmsg/agentmsg_grpc.pb.go | bin
	go build -v -buildvcs=false -ldflags="${LDFLAGS}" -o bin/agent ./cmd/agent

bin:
	mkdir -v bin

pkg/agentmsg:
	mkdir -pv pkg/agentmsg

pkg/agentmsg/agentmsg.pb.go pkg/agentmsg/agentmsg_grpc.pb.go: agentmsg.proto pkg/agentmsg
	protoc -I. --go_out=./pkg/agentmsg --go_opt=paths=source_relative --go-grpc_out=./pkg/agentmsg --go-grpc_opt=paths=source_relative agentmsg.proto