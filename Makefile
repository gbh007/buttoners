BUILD_ENV = GOOS=linux GOARCH=amd64 CGO_ENABLED=0
WD = $(shell pwd)
SERVICE_LIST = $(shell ls -d services/*/)
GO_APP_LIST = $(SERVICE_LIST) core/ ui/console

.PHONY: gen-proto
gen-proto:
	protoc -I=. --go_out=. --go-grpc_out=. services/gate/gate.proto

.PHONY: generate
generate: gen-proto

.PHONY: install-proto
install-proto:
	sudo apt install protobuf-compiler
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@lates

.PHONY: build
build:
	$(BUILD_ENV) go build -o ./bin/build/gate ./services/gate/cmd/gate
	$(BUILD_ENV) go build -o ./bin/build/auth ./services/auth/cmd/authserver
	$(BUILD_ENV) go build -o ./bin/build/handler ./services/handler/cmd/handler
	$(BUILD_ENV) go build -o ./bin/build/worker ./services/worker/cmd/worker
	$(BUILD_ENV) go build -o ./bin/build/log ./services/log/cmd/log
	$(BUILD_ENV) go build -o ./bin/build/notification ./services/notification/cmd/notification
	$(BUILD_ENV) go build -o ./bin/build/console-ui ./ui/console/cmd

.PHONY: up
up: build
	docker compose -f ./deployments/docker-compose.yml up -d --build --remove-orphans

.PHONY: logs
logs:
	docker compose -f ./deployments/docker-compose.yml logs -f auth gate handler log notification worker

.PHONY: down
down:
	docker compose -f ./deployments/docker-compose.yml down --remove-orphans

.PHONY: cleanup
cleanup:
	docker compose -f ./deployments/docker-compose.yml down --remove-orphans --volumes
	rm -r ./bin/build

.PHONY: console-ui
console-ui: build
	./bin/build/console-ui

.PHONY: lint
lint:
	golangci-lint run


.PHONY: tidy
tidy:
	for path in $(GO_APP_LIST); do \
		cd $(WD)/$$path; \
		echo $$path; \
		go mod tidy; \
	done