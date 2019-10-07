# Work at Olist

This project aims to develop a Rest API to control phone call records and generate bills.

# Getting Started

These instructions will provide you a copy of the project that can be run on your local machine for development and testing purposes.
Consult [deployment](#deployment) item for notes on how to deploy the project on a live system.

# Prerequisites

This package is built with go1.13, and all you need is provide with the go standard library.

# Installing

This is what you need to install the application from the source code:

```shell script
    git clone https://github.com/paraizofelipe/work-at-olist.git
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

# API

## POST - /api/records

### Body of start requests

```json
{
	"type": "start",
	"timestamp": "2019-10-06T05:00:00Z",
	"call_id": 1,
	"source": "4199999999",
	"destination": "4188888888"
}
```

### Body of end requests

```json
{
	"type": "end",
	"timestamp": "2019-10-06T06:00:00Z",
	"call_id": 1
}
```

### Body of response

**Status**: 201

```json
{
    "message": "record successfully saved"
}
```

**Status**: 422

```json
{
    "error": "..."
}
```

## GET - /api/bills/

#### Path params

/api/bills/\<subscriber\>

| param      | type   | example     |
|------------|--------|-------------|
| subscriber | number | 41991233821 |

#### Query string

| query | type   | example |
|-------|--------|---------|
| year  | number | 2019    |
| month | number | 05      |

### Body of response

**Status**: 200

```json
{
  "subscriber": "99988526423",
  "month": 2,
  "year": 2016,
  "calls": [
    {
      "bill_id": 1,
      "destination": "9933468278",
      "duration": "2h0m0s",
      "start_date": "2016-02-29",
      "start_time": "12:0:0",
      "price": 11.16
    }
  ],
  "price": 11.16
}
```
 **Status**: 404
 
 ```json
{
     "errors": {
         "message": "bills not found"
     }
}
 ```

# Built With

* [Arch linux](https://www.archlinux.org/) - Operational sistem
* [Neovim](https://neovim.io/)  - Editor Text
* [go](https://golang.org/) - The GO programing language.
* [go-sqlite3](https://github.com/mattn/go-sqlite3) - Driver for sqlite3 database 

# Versioning

This project use SemVer for versioning. For the versions available, see the tags on this repository.

# Authors

Felipe Paraizo - Initial work - [paraizo](http://paraizo.dev)

