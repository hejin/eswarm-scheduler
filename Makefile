all: local

local: fmt vet
	GOOS=linux GOARCH=amd64 go build  -o=bin/eswarm-scheduler ./cmd/scheduler

build:  local
	docker build --no-cache . -t registry.cn-hangzhou.aliyuncs.com/oprobot/eswarm-scheduler:0.0.1

push:   build
	docker push registry.cn-hangzhou.aliyuncs.com/oprobot/eswarm-scheduler:0.0.1

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean: fmt vet
	sudo rm -f eswarm-scheduler
