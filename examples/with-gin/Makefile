NAME=supertokens-golang-demo
VERSION=0.0.1
BUILD=build

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME) -tags=jsoniter -o $(BUILD)/$(NAME)

.PHONY: run
## run: Build and Run in development mode.
run: build
	@./$(BUILD)/$(NAME) -e dev

.PHONY: run-prod
## run-prod: Build and Run in production mode.
run-prod: build
	@./$(BUILD)/$(NAME) -e prod

.PHONY: run-stage
## run-stage: Build and Run in staging mode.
run-stage: build
	@./$(BUILD)/$(NAME) -e stage

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(BUILD)/$(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

# .PHONY: test
# ## test: Run tests with verbose mode
# test:
# 	@go test -v ./tests/*

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
