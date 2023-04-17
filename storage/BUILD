load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "user",
    srcs = ["user.go"],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//auth:password",
        "//proto:user_go_proto",
        "@com_github_google_uuid//:uuid",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "storage_manager",
    srcs = [
        "storage_manager.go",
        "user.go",
    ],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//auth:password",
        "//proto:user_go_proto",
        "@com_github_google_uuid//:uuid",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_test(
    name = "storage_manager_test",
    srcs = ["storage_manager_test.go"],
    embed = [
        ":storage_manager",
    ],
    deps = [
        "//proto:user_go_proto",
        "@com_github_proullon_ramsql//driver",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
        "@io_gorm_driver_postgres//:postgres",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
