# Work at Olist

This project aims to develop a Rest API to control phone call records and generate bills.

# Getting Started

These instructions will provide you a copy of the project that can be run on your local machine for development and testing purposes.
Consult deployment item for notes on how to deploy the project on a live system.

# Prerequisites

This package is built with go1.13, and all you need is provided with the go standard library.

# Installing

This is what you need to install the application from the source code:

```shell script
    go install
```

To build the docker version you can use the `Makefile`:

```shell script
    make dk-build 
```

# Running the tests

Until I finish this README there is not so much Unit tests written.

But I will try to coverage unless 70% of unit tests for this code as soon as possible.

You can run tests like this:

```shell script
    make test
```

To run this code locally for test purposes use:

```shell script
    PORT=8989 DEBUG=true go run main.go
```

# Deployment

This codebase is cloud-native by design so you can use lots of environments to make this run anywhere you want.

But to make this even easier to you the codebase also provides a Dockerfile and a docker-compose.

There is also a Makefile to make all this even easier.

Deploy with docker-compose:
```shell script
    docker-compose up --build
```

Run with Makefile:

```shell script
    make run
```

# Built With

[go](https://golang.org/) - The GO programing language.
[go-sqlite3](https://github.com/mattn/go-sqlite3) - driver for sqlite3 database 

# Versioning

This project use SemVer for versioning. For the versions available, see the tags on this repository.

# Authors

Felipe Paraizo - Initial work - [paraizo](http://paraizo.dev)
