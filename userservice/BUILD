load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "userservice",
    srcs = ["userservice.go"],
    importpath = "github.com/mnadev/limestone/userservice",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:user_service_go_proto",
        "//proto:user_go_proto",
        "//storage:storage_manager",
    ],
)
