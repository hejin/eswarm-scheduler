all: local

local: fmt vet
	GOOS=linux GOARCH=arm64 go build  -o=bin/eswarm-scheduler ./cmd/scheduler

build:  local
	docker build --no-cache . -t o5jeff/eswarm-scheduler:0.0.3
#	docker build --no-cache . -t registry.cn-hangzhou.aliyuncs.com/oprobot/eswarm-scheduler:0.0.3

push:   build
	docker push o5jeff/eswarm-scheduler:0.0.3
#	docker push registry.cn-hangzhou.aliyuncs.com/oprobot/eswarm-scheduler:0.0.3

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean: fmt vet
	sudo rm -f eswarm-scheduler
