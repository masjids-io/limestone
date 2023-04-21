load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "user_service",
    srcs = ["user_service.go"],
    importpath = "github.com/mnadev/limestone/user_service",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:user_service_go_proto",
        "//storage:storage_manager",
    ],
)
