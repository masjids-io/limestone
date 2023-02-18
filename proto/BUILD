load("@rules_proto_grpc//go:defs.bzl", "go_grpc_library", "go_proto_library")

proto_library(
    name = "auth_service_proto",
    srcs = ["auth_service.proto"],
    deps = [":user_proto"],
)

proto_library(
    name = "user_proto",
    srcs = ["user.proto"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
    ],
)

proto_library(
    name = "user_service_proto",
    srcs = ["user_service.proto"],
    deps = [
        ":user_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

go_grpc_library(
    name = "auth_service_go_proto",
    importpath = "github.com/mnadev/limestone",
    protos = [":auth_service_proto"],
)

go_proto_library(
    name = "user_go_proto",
    importpath = "github.com/mnadev/limestone",
    protos = [":user_proto"],
)

go_grpc_library(
    name = "user_service_go_proto",
    importpath = "github.com/mnadev/limestone",
    protos = [
        ":user_proto",
        ":user_service_proto",
    ],
)
