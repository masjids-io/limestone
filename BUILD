load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/mnadev/limestone
gazelle(name = "gazelle")

gazelle(
    name = "gazelle_update_repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

go_library(
    name = "limestone_lib",
    srcs = ["main.go"],
    importpath = "github.com/mnadev/limestone",
    visibility = ["//visibility:private"],
    deps = ["@org_golang_google_grpc//:go_default_library"],
)

go_binary(
    name = "limestone",
    embed = [":limestone_lib"],
    visibility = ["//visibility:public"],
)
