load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_proto_grpc//go:defs.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "auth_service_proto",
    srcs = ["auth_service.proto"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
    ],
)

proto_library(
    name = "user_service_proto",
    srcs = ["user_service.proto"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "auth_service_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/mnadev/limestone/authservice/proto",
    protos = [":auth_service_proto"],
)

go_grpc_library(
    name = "user_service_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/mnadev/limestone/user_service/proto",
    protos = [
        ":user_service_proto",
    ],
    deps = [],
)
