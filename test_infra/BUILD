load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "unit_test_suite",
    testonly = 1,
    srcs = ["unit_test_suite.go"],
    importpath = "github.com/mnadev/limestone/test_infra",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:user_service_go_proto",
        "//storage:storage_manager",
        "//user_service:user_service",
        "@com_github_proullon_ramsql//driver",
        "@com_github_stretchr_testify//suite",
        "@io_gorm_driver_postgres//:postgres",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//test/bufconn",
    ],
)

go_library(
    name = "integration_test_suite",
    testonly = 1,
    srcs = ["integration_test_suite.go"],
    importpath = "github.com/mnadev/limestone/test_infra",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:user_service_go_proto",
        "//storage:storage_manager",
        "//user_service:user_service",
        "@com_github_proullon_ramsql//driver",
        "@com_github_stretchr_testify//suite",
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
        ":integration_test_suite",
    ],
    deps = [
        "//proto:user_service_go_proto",
        "//user_service:user_service",
        "@com_github_proullon_ramsql//driver",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//suite",
        "@io_gorm_driver_postgres//:postgres",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_test(
    name = "storage_manager_test",
    srcs = ["storage_manager_test.go"],
    embed = [
        ":unit_test_suite",
    ],
    deps = [
        "//proto:user_service_go_proto",
        "//storage:storage_manager",
        "@com_github_proullon_ramsql//driver",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//suite",
        "@io_gorm_driver_postgres//:postgres",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
