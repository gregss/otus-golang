.PHONY: up
up:
	docker-compose -f docker-compose.yml up -d --build

.PHONY: down
down:
	docker-compose down

.PHONY: build-app
build-app:
	CGO_ENABLED=0 go build -o ./bin/ ./cmd/app/app.go

.PHONY: build-scheduler
build-scheduler:
	CGO_ENABLED=0 go build -o ./bin/ ./cmd/scheduler/scheduler.go

.PHONY: build-sender
build-sender:
	CGO_ENABLED=0 go build -o ./bin/ ./cmd/sender/sender.go

.PHONY: build
build: build-app build-scheduler build-sender

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: fix-lint
fix-lint:
	gci write .
	gofumpt -l -w .

generate:
	protoc \
		--proto_path=api/ \
		--go_out=api/ \
		--go-grpc_out=api/ \
		api/*.proto

.PHONY: build run build-img run-img version test lint

.PHONY: integration-testsn
integration-tests:
	go test -tags integrations ./tests/...

# вывод интеграционных тестов
# https://i.imgur.com/sgbOQ0Q.png
# https://i.imgur.com/4K1f7hY.png
.PHONY: up-test
up-test:
	docker-compose -f docker-compose.test.yml up -d --build