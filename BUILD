load("@io_bazel_rules_docker//go:image.bzl", "go_image")
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
    name = "server_lib",
    srcs = ["server.go"],
    importpath = "github.com/mnadev/limestone",
    visibility = ["//visibility:private"],
    deps = [
        "//event_service:event_service",
        "//masjid_service:masjid_service",
        "//proto:event_service_go_proto",
        "//proto:masjid_service_go_proto",
        "//proto:user_service_go_proto",
        "//storage:storage_manager",
        "//user_service:user_service",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@io_gorm_driver_postgres//:postgres",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "server",
    embed = [":server_lib"],
    visibility = ["//visibility:public"],
)

go_image(
    name = "server_image",
    binary = ":server",
)
