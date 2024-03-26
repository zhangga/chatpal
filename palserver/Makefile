# 命令
GO_BUILD = CGO_ENABLED=0 go build -trimpath
GO_RUN = CGO_ENABLED=0 go run -trimpath -ldflags "-s -w"

MODEL_PROTO_FILES=$(shell find model -name *.proto)

.PHONY: test
# run all test
test:
	@make check-package
	go test -race ./...

.PHONY: lint
# run all lint
lint:
	@make check-package
	golangci-lint run -c .golangci.yml ./...

.PHONY: model
model:
	protoc 	--proto_path=./model \
			--proto_path=../third_party \
		  	--gogo_out=paths=source_relative:./model \
			$(MODEL_PROTO_FILES)

.PHONY: example
example:
	$(GO_RUN) example/main.go

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
