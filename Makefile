# --- Variables ---

# Include .env file if available:
-include .env

# Fake AWS credentials as fix for AWS SAM Local issue #134:
# See also https://github.com/awslabs/aws-sam-local/issues/134
FAKE_AWS_ENV = AWS_ACCESS_KEY_ID=0 AWS_SECRET_ACCESS_KEY=0

# AWS CLI wrapped with aws-vault for secure credentials access,
# can be overriden by defining the AWS_CLI environment variable:
AWS_CLI ?= aws-vault exec '$(AWS_PROFILE)' -- aws

# The import path of the passphrase package:
IMPORT_PATH = github.com/blueimp/passphrase

# The absolute package path:
PKG_PATH = $(GOPATH)/src/$(IMPORT_PATH)

# The absolute path for the passphrase binary installation:
BIN_PATH = $(GOPATH)/bin/passphrase

# Dependencies to build the passphrase command-line interface:
CLI_DEPS = $(PKG_PATH) passphrase/cli.go passphrase.go words.go

# Dependencies to build the lambda application:
LAMBDA_DEPS = $(PKG_PATH) vendor lambda/lambda.go passphrase.go words.go


# --- Main targets ---

# The default target builds both the CLI and lambda binaries:
all: passphrase lambda

# Builds the CLI binary:
passphrase: passphrase/passphrase

# Cross-compiles the lambda binary:
lambda: lambda/bin/main

# Generates the word list as go code:
words:
	go generate

# Runs the unit tests for all components:
test: words.go
	@go test .
	@cd passphrase; go test .
	@cd lambda; go test .

# Installs the passphrase binary at $GOPATH/bin/passphrase:
install: $(BIN_PATH)

# Deletes the passphrase binary from $GOPATH/bin/passphrase:
uninstall:
	rm -f $(BIN_PATH)

# Generates a sample lambda event:
event: lambda/event.json

# Invokes the lambda function locally:
invoke: lambda/event.json lambda/bin/main
	cd lambda; $(FAKE_AWS_ENV) sam local invoke -e event.json

# Starts a local API Gateway for the lambda function as background process:
start: lambda/bin/main
	cd lambda; $(FAKE_AWS_ENV) sam local start-api & echo "$$!" > server.pid

# Stops the local API Gateway:
stop:
	cd lambda; test -f server.pid && kill "$$(cat server.pid)" && rm server.pid

# Starts the local API Gateway and a watch process for the lambda function,
# on MacOS also automatically reloads the active Chrome/Safari/Firefox tab:
watch:
	@exec ./watch.sh $(BROWSER)

# Deploys the lambda function to AWS:
deploy: lambda/deployed.txt url

# Prints the API Gateway URL of the deployed lambda function:
url: lambda/passphrase.url
	@grep -o 'https://.*' lambda/passphrase.url

# Deletes the CloudFormation stack of the lambda function:
destroy:
	rm -f lambda/passphrase.url
	$(AWS_CLI) cloudformation delete-stack --stack-name '$(STACK_NAME)'

# Removes all build artifacts:
clean:
	rm -f \
		lambda/bin/main \
		lambda/debug \
		lambda/debug.test \
		lambda/deploy.yml \
		lambda/deployed.txt \
		lambda/event.json \
		lambda/passphrase.url \
		passphrase/debug \
		passphrase/debug.test \
		passphrase/passphrase


# --- Helper targets ---

# Defines phony targets (targets without a corresponding target file):
.PHONY: \
	all \
	passphrase \
	lambda \
	words \
	test \
	install \
	uninstall \
	event \
	invoke \
	start \
	stop \
	watch \
	deploy \
	url \
	destroy \
	clean

# Creates a symlink from the project import path to this directory.
# This allows working with a project outside of $GOPATH:
$(PKG_PATH):
	mkdir -p $@; cd $@/..; rm -rf passphrase; ln -s '$(PWD)'

# Installs the passphrase binary at $GOPATH/bin/passphrase:
$(BIN_PATH): $(CLI_DEPS)
	go install $(IMPORT_PATH)/passphrase

# Install dependencies via `dep ensure` if available, else via `go get`:
vendor: $(PKG_PATH)
	if command -v dep > /dev/null 2>&1; then cd $(PKG_PATH); dep ensure; \
		else go get -d ./... && mkdir vendor; fi

# Builds the CLI binary:
passphrase/passphrase: $(CLI_DEPS)
	cd passphrase; go build

# Cross-compiles the lambda binary:
# ldflags explanation (see `go tool link`):
#   -s  disable symbol table
#   -w  disable DWARF generation
lambda/bin/main: $(LAMBDA_DEPS)
	cd $(PKG_PATH)/lambda; \
		GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/main

# Generates the word list as go code if generate.go or words.txt change:
words.go: generate.go words.txt
	go generate

# Generates a sample lambda event:
lambda/event.json:
	cd lambda; sam local generate-event api > event.json

# Packages the lambda binary and uploads it to S3:
lambda/deploy.yml: lambda/bin/main lambda/template.yml
	cd lambda; $(AWS_CLI) cloudformation package \
		--template-file template.yml \
		--s3-bucket '$(DEPLOYMENT_BUCKET)' \
		--s3-prefix '$(DEPLOYMENT_PREFIX)' \
		--output-template-file deploy.yml

# Deploys the packaged binary to AWS:
lambda/deployed.txt: lambda/deploy.yml
	cd lambda; $(AWS_CLI) cloudformation deploy \
		--template-file deploy.yml \
		--stack-name '$(STACK_NAME)' \
		--parameter-overrides LambdaRole='$(LAMBDA_ROLE)'
	date >> lambda/deployed.txt

# Generates a passphrase.url file with the API Gateway URL:
lambda/passphrase.url:
	API_GW_ID=$$($(AWS_CLI) cloudformation describe-stack-resource \
		--stack-name '$(STACK_NAME)' \
		--logical-resource-id ServerlessRestApi \
		--query StackResourceDetail.PhysicalResourceId \
		--output text \
	) && \
	printf '%s\nURL=https://%s.execute-api.$(AWS_REGION).amazonaws.com/Prod\n' \
		[InternetShortcut] \
		"$$API_GW_ID" \
		> lambda/passphrase.url
