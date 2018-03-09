# Passphrase
> Better passwords by combining random words.

Passphrase is a Go library, command-line interface and AWS Lambda function to
generate a random sequence of words.

It is inspired by Randall Munroe's [xkcd webcomic #936](https://xkcd.com/936/)
with the title "Password Strength":

![Password Strength](https://imgs.xkcd.com/comics/password_strength.png)

## Installation
The `passphrase` command-line interface can be installed via
[go get](https://golang.org/cmd/go/):

```sh
go get github.com/blueimp/passphrase/passphrase
```

## Usage
By default, `passphrase` prints four space-separated words, but also accepts
an argument for the number of words to generate:

```sh
passphrase [number]
```

## Word list
This repository includes the word list `google-10000-english-usa-no-swears.txt`
from Josh Kaufman's repository
[google-10000-english](https://github.com/first20hours/google-10000-english/),
but `passphrase` can also be compiled with another list of newline separated
words.

The words module can be generated the following way:

```sh
WORD_LIST_URL=words.txt MIN_WORD_LENGTH=3 make words
```

The `WORD_LIST_URL` variable can point to a URL or a local file path and
falls back to `words.txt`.

Words shorter than `MIN_WORD_LENGTH` (defaults to a minimum word
length of `3` characters) are filtered out.

The updated word list module can then be used in a new build.

## Build
To build both the CLI and the AWS Lambda function binary, run
[Make](https://en.wikipedia.org/wiki/Make_\(software\)) in the repository:

```sh
make
```

Both components can also be built individually:

```sh
make passphrase
```

```sh
make lambda
```

The locally built binary can be installed at `$GOPATH/bin/passphrase` with the
following command:

```sh
make install
```

The uninstall command removes the binary from `$GOPATH/bin/passphrase`:

```sh
make uninstall
```

To clean up all build artifacts, run the following:

```sh
make clean
```

## Test
All components come with unit tests, which can be executed the following way:

```sh
make test
```

## AWS Lambda
Passphrase can be deployed as [AWS lambda](https://aws.amazon.com/lambda/)
function with an [API Gateway](https://aws.amazon.com/api-gateway/) triggger.

The function accepts a query parameter `n` to define the number of words to
generate, but limits the sequence to `100` words, e.g.:

```
https://API_GW_ID.execute-api.REGION.amazonaws.com/Prod?n=100
```

### Requirements
Deployment requires the [AWS CLI](https://aws.amazon.com/cli/) as well as
[aws-vault](https://github.com/99designs/aws-vault) for secure credentials
access.  
Alternatively, it's also possible to reset the wrapped `aws` CLI command by
exporting `AWS_CLI=aws` as environment variable.

Local invocations require
[AWS SAM Local](https://github.com/awslabs/aws-sam-local).

The local watch task requires [entr](https://bitbucket.org/eradman/entr) to be
installed, which is available in the repositories of popular Linux distributions
and can be installed on MacOS via [Homebrew](https://brew.sh/):

```sh
brew install entr
```

### Environment variables
To be able to deploy, the following variables have to be set, e.g. by adding
them to a `.env` file, which gets included in the provided `Makefile`:

```sh
# The AWS profile to use for aws-vault:
AWS_PROFILE=default
# The S3 bucket where the lambda package can be uploaded:
DEPLOYMENT_BUCKET=example-bucket
# The S3 object prefix for the lambda package:
DEPLOYMENT_PREFIX=passphrase
# The CloudFormation stack name:
STACK_NAME=passphrase
# The name of an existing IAM role for AWS Lambda with
# AWSLambdaBasicExecutionRole attached:
LAMBDA_ROLE=arn:aws:iam::000000000000:role/aws-lambda-basic-execution-role
# The AWS service region, required to construct the API Gateway URL:
AWS_REGION=eu-west-1
```

### Deploy
To package and deploy the function, execute the following:

```sh
make deploy
```

After the deployment succeeds, the [API Gateway](https://aws.amazon.com/api-gateway/) URL is printed.

This URL can also be retrieved later with the following command:

```sh
make url
```

To remove the AWS Lambda function and API Gateway configuration, execute the
following:

```sh
make destroy
```

### Local development
Using [AWS SAM Local](https://github.com/awslabs/aws-sam-local), the function
can also be invoked and served locally.

A sample API Gateway event can be generated the following way:

```sh
make event
```

To invoke the function locally, execute the following:

```sh
make invoke
```

To start a local API Gateway in a background process, use the following command:

```sh
make start
```

To stop the local API Gateway background process again, execute this command:

```sh
make stop
```

To start the local API Gateway along with a watch process for source file
changes, run the following:

```sh
[BROWSER=chrome|safari|firefox] make watch
```

The watch task recompiles the lambda binary on changes.  
On MacOS, it also automatically reloads the active Chrome/Safari/Firefox tab.

## License
Released under the [MIT license](https://opensource.org/licenses/MIT).
