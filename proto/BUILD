load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_proto_grpc//grpc-gateway:defs.bzl", "gateway_grpc_library", "gateway_openapiv2_compile")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "adhan_service_proto",
    srcs = ["adhan_service.proto"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
        "@go_googleapis//google/api:annotations_proto",
    ],
)

proto_library(
    name = "event_service_proto",
    srcs = ["event_service.proto"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
        "@go_googleapis//google/api:annotations_proto",
    ],
)

proto_library(
    name = "masjid_service_proto",
    srcs = ["masjid_service.proto"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
        "@go_googleapis//google/api:annotations_proto",
    ],
)

proto_library(
    name = "user_service_proto",
    srcs = ["user_service.proto"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
        "@go_googleapis//google/api:annotations_proto",
    ],
)

gateway_grpc_library(
    name = "adhan_service_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/mnadev/limestone/adhan_service/proto",
    protos = [":adhan_service_proto"],
)

gateway_grpc_library(
    name = "event_service_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/mnadev/limestone/event_service/proto",
    protos = [":event_service_proto"],
)

gateway_grpc_library(
    name = "masjid_service_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/mnadev/limestone/masjid_service/proto",
    protos = [":masjid_service_proto"],
)

gateway_grpc_library(
    name = "user_service_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/mnadev/limestone/user_service/proto",
    protos = [":user_service_proto"],
)

gateway_openapiv2_compile(
    name = "api_openapiv2",
    protos = [
        ":adhan_service_proto",
        ":event_service_proto",
        ":masjid_service_proto",
        ":user_service_proto",
    ],
)
