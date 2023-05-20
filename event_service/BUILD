load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "event_service",
    srcs = ["event_service.go"],
    importpath = "github.com/mnadev/limestone/event_service",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:event_service_go_proto",
        "//storage:storage_manager",
    ],
)
