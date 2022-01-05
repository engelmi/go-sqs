# go-sqs

A simplistic library to listen for and produce messages to Amazons **S**imple **Q**ueuing **S**ervice (SQS). 

![alt text](https://imgur.com/VbbjtZ7.png "Messaging Gopher")

![Build Status](https://travis-ci.com/engelmi/go-sqs.svg?branch=main)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## Quick start
Get started with `go-sqs` via
```bash
$ go get github.com/engelmi/go-sqs
```

## Features
* Provides an easy to use producer interface to put messages into an SQS queue
  * Can be configured via a configuration struct
* Provides an easy to use consumer interface to listen for messages in an SQS queue
  * Can be configured via a configuration struct, e.g. with poll timeout
  * A handler function is used to process messages
  * Errors returned from the handler function - as well as any occuring panic - result in the message **not** being acknowledged
  * Messages that are not acknowledged are not deleted from the queue and, therefore, consumed again
  * The listening for messages can be started as a goroutine
  * A wait group can be used to have the currently executed handler finish before terminating the program


## Usage
Examples on how to use `go-sqs` are provided in `github.com/engelmi/go-sqs/examples`. 

The credentials file `~/.aws/credentials` containing dummy values for the access key and secret key is required to run the examples: 
```
[default]
aws_access_key_id = foo
aws_secret_access_key = bar
```
Alternatively, provide both values via environment variables: 
```bash
AWS_ACCESS_KEY_ID='foo' AWS_SECRET_ACCESS_KEY='bar' go run <example-file>.go
```
