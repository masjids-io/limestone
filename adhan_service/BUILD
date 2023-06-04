load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "adhan_service",
    srcs = ["adhan_service.go"],
    importpath = "github.com/mnadev/limestone/adhan_service",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:adhan_service_go_proto",
        "//storage:storage_manager",
    ],
)
