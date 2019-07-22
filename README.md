# Bitcoin Address Generator

This service is intended to be run with AWS Lambda and DynamoDB.  
Check out this [tutorial](https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda) to set up the AWS environment in order to have the Lambda function allowed to work along with a DynamoDB table.

## Requirements

* Golang
* Dep

## Installation

Install dependencies
```sh
$ bash scripts/install
```

Build the binary

```sh
$ bash scripts/build linux
```

Add the binary to a zip archive to prepare it for the uploading on AWS Lambda

```sh
zip -j build/function.zip build/btc-address-generator-linux
```


## Usage

The following environment variables must be set:

* `XPUB` the extended public key to derive address from
* `NETWORK` the bitcoin network (mainnet, testnet or regtest)
* `DERIVATON_PATH` the *change* and/or *account index* of the bip32 derivation path (e.g. `m/0/0` or `m/0`)

Example:

```sh
$ XPUB=upub5FxWB33rLzN7YHBFqbu9tB9Xx26T2QZPHExHE9zqVYnoBQpryFqFMMZPMh7hezCMahwDfHcowFhvNMXab494VqgUiby5Z6xQ3c8ZpxgSbwt NETWORK=testnet DERIVATION_PATH=m/0 ./build/btc-address-generator-linux
```

You must create a table `BtcAddressGenerator` within DynamoDB as the Lambda expects it's already created.  
The function automatically checks if an item in the table already exists for the current extended key, otherwise it creates a new one before returning the very first generated address.

## Test

```sh
$ go test -v
```

Tests do not provide a dry run, thus for every test case a new item is created into the database. You have to remove the items manually to keep things clean.

If you have more than one AWS profile specified in the `~/.aws/credentials` file, run tests passing the `AWS_PROFILE` environment variable:

```sh
AWS_PROFILE=altafan go test -v
```