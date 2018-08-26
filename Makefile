# --- Variables ---

# Include .env file if available:
-include .env

# The platform to use for local development and deployment.
# Can be either "appengine" or "lambda":
PLATFORM ?= appengine

# Fake AWS credentials as fix for AWS SAM Local issue #134:
# See also https://github.com/awslabs/aws-sam-local/issues/134
FAKE_AWS_ENV = AWS_ACCESS_KEY_ID=0 AWS_SECRET_ACCESS_KEY=0

# AWS CLI wrapped with aws-vault for secure credentials access,
# can be overriden by defining the AWS_CLI environment variable:
AWS_CLI ?= aws-vault exec '$(AWS_PROFILE)' -- aws

# The absolute path for the passphrase binary installation:
BIN_PATH = $(GOPATH)/bin/passphrase

# Dependencies to build the passphrase command-line interface:
CLI_DEPS = passphrase/cli.go passphrase/go.mod passphrase.go words.go

# Dependencies to build the lambda application:
LAMBDA_DEPS = lambda/lambda.go lambda/go.mod passphrase.go words.go


# --- Main targets ---

# The default target builds the CLI binary:
all: passphrase/passphrase

# Cross-compiles the lambda binary:
lambda: lambda/bin/main

# Generates the word list as go code:
words:
	go generate

# Runs the unit tests for all components:
test: words.go
	@go test ./...
	@cd passphrase; go test ./...
	@cd appengine; go test ./...
	@cd lambda; go test ./...

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

# Starts a local server for the given platform:
start: $(PLATFORM)-start

# Starts a local server for the given platform and a watch process:
watch: $(PLATFORM)-watch

# Deploys the project for the given platform:
deploy: $(PLATFORM)-deploy

# Opens a browser tab with the production URL of the App Engine project:
browse:
	cd appengine; gcloud app browse --project $(PROJECT)

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
	watch \
	deploy \
	browse \
	url \
	destroy \
	appengine-start \
	appengine-watch \
	appengine-deploy \
	lambda-start \
	lambda-watch \
	lambda-deploy \
	clean

# Installs the passphrase binary at $GOPATH/bin/passphrase:
$(BIN_PATH): $(CLI_DEPS)
	cd passphrase; go install

# Builds the passphrase binary:
passphrase/passphrase: $(CLI_DEPS)
	cd passphrase; go build

# Generates the word list as go code if generate.go or words.txt change:
words.go: generate.go words.txt
	go generate

# Starts a local App Engine server:
appengine-start:
	cd appengine; dev_appserver.py .

# Starts a local App Engine server and a watch process for source file changes,
# on MacOS also automatically reloads the active Chrome/Safari/Firefox tab:
appengine-watch:
	@exec ./watch.sh start $(BROWSER)

# Deploys the App Engine project to Google Cloud:
appengine-deploy:
	cd appengine; gcloud app deploy --project $(PROJECT) --version $(VERSION)

# Starts a local API Gateway:
# Fake AWS credentials as fix for AWS SAM Local issue #134:
# See also https://github.com/awslabs/aws-sam-local/issues/134
lambda-start:
	cd lambda; AWS_ACCESS_KEY_ID=0 AWS_SECRET_ACCESS_KEY=0 sam local start-api

# Starts a local API Gateway and a watch process for source file changes,
# on MacOS also automatically reloads the active Chrome/Safari/Firefox tab:
lambda-watch:
	@exec ./watch.sh start $(BROWSER) -- make -s lambda

# Deploys the lambda function to AWS:
lambda-deploy: lambda/deployed.txt url

# Cross-compiles the lambda binary:
# ldflags explanation (see `go tool link`):
#   -s  disable symbol table
#   -w  disable DWARF generation
lambda/bin/main: $(LAMBDA_DEPS)
	cd lambda; \
		GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o bin/main

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
