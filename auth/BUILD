load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "password",
    srcs = ["password.go"],
    importpath = "github.com/mnadev/limestone/auth",
    visibility = ["//visibility:public"],
    deps = ["@org_golang_x_crypto//bcrypt"],
)

go_test(
    name = "password_test",
    srcs = ["password_test.go"],
    embed = [":password"],
    deps = [
        "@com_github_stretchr_testify//assert",
    ],
)
