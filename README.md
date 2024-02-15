# limestone
Backend server for masjids.io

## masjids.io
masjids.io is an all-in-one masjid technological, administration and management tool that enables masjids to easily establish an online footprint, and manage the services they offer through app integrations.

The current goal of masjids.io is to create a masjid administration tool, website management and hosting tool and a companion adhan app.

Limestone is the backend server for masjids.io which enables all of this functionality. It has (or will have) a user service
with OAuth logins and SSO support, masjid registration, setting roles for different users for a masjid, role based authorization, website hosting and much more. 

## Prerequisites

In order to run the server, you will need to install and be familiar with the following:

* [Go](https://go.dev/)
* [Protobuf](https://protobuf.dev/downloads/)
* [Buf](https://buf.build/)
* [Docker](https://bazel.build/)

## Regenerating protos

To regenerate protos, run `cd proto`, then run `buf generate .`. This MUST be done after each proto change.

## Regenerate swagger docs

To regenerate the Swagger YAML, run `cd proto`, then run `buf generate .`. This MUST be done after each proto change.

To create the `index.html`` file from the YAML file, first install redoc-cli via:

`npm install -g redoc-cli`

Then, run:

`redoc-cli bundle -o docs/index.html docs/apidocs.swagger.yaml`


## Running the server
To run the server, make sure the relevant environment variables are set; then, run the command:

`go run server.go`

from the root directory. This exposes both an HTTP server and a gRPC server. You can make calls to the gRPC server via [grpcurl](https://github.com/fullstorydev/grpcurl). But, for end-to-end testing, it's just easier to call the HTTP server, with any HTTP client (Postman or curl).

## Running dockerized environment
To run the server with the databases, make sure the relevant environment variables are set; and you have Make installed along with Docker run the command:

To rebuild binaries and start containers:
`make up_build`

To only start pre-built binaries:
`make up`

To bring down the docker environment:
`make down`

## Tasks

### Implemented

- User service
- Event service
- Masjid Service
- Adhan service
- gRPC Integration tests

### TODOs

- Nikkah Service API
- Reverts Service API
- JWT authentication
- HTTP integration tests
- Docker-ify server
- Setup github workflows to test and check for formatting, and push to prod
- Host container on Cloud
- Authorization service
- OAuth 2 SSO for providers like Google, Facebook, Microsoft, Apple.
- Custom website builder service and rendering server
- Adhan livestreaming API
