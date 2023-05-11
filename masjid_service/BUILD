load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "masjid_service",
    srcs = ["masjid_service.go"],
    importpath = "github.com/mnadev/limestone/masjid_service",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:masjid_service_go_proto",
        "//storage:storage_manager",
    ],
)
