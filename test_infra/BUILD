load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "test_server",
    srcs = ["test_server.go"],
    testonly = 1,
    importpath = "github.com/mnadev/limestone/test_infra",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:user_service_go_proto",
        "//storage:storage_manager",
        "//user_service:user_service",
        "@com_github_proullon_ramsql//driver",
        "@io_gorm_driver_postgres//:postgres",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//test/bufconn",
    ],
)

go_test(
    name = "user_service_test",
    srcs = ["user_service_test.go"],
    embed = [
        ":test_server",
    ],
    deps = [
        "//proto:user_service_go_proto",
        # "//storage:storage_manager",
        "//user_service:user_service",
        "@com_github_proullon_ramsql//driver",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
        "@io_gorm_driver_postgres//:postgres",
        # "@io_gorm_gorm//:gorm",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
