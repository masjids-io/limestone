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

NOTE: I may remove Bazel from this repository and move to using Buf to generate the gRPC code.

## Regenerating protos

To regenerate the Go proto libraries, first run `cd proto` to go into the proto directory. Then run `buf generate .` to regenerate the proto libraries. You should do this every time you make a change to the protos.

## Running the server
To run the server, run the command:

`bazel run //:server`

## Docker
To build the container image, run:

`bazel build //:server_image --@io_bazel_rules_docker//transitions:enable=false` 

To run the container image, run: 

`bazel run //:server_image --@io_bazel_rules_docker//transitions:enable=false` 
