load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "user",
    srcs = ["user.go"],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:user_go_proto",
        "@com_github_google_uuid//:uuid",
        "@io_gorm_gorm//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_x_crypto//bcrypt:go_default_library",
    ],
)
