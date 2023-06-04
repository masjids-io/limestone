load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

# Libraries

go_library(
    name = "assert",
    testonly = 1,
    srcs = ["assert.go"],
    importpath = "github.com/mnadev/limestone/test_infra",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:adhan_service_go_proto",
        "//proto:event_service_go_proto",
        "//proto:masjid_service_go_proto",
        "//proto:user_service_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@com_github_stretchr_testify//assert",
        "@org_golang_google_protobuf//testing/protocmp",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "integration_test_suite",
    testonly = 1,
    srcs = ["integration_test_suite.go"],
    importpath = "github.com/mnadev/limestone/test_infra",
    visibility = ["//visibility:public"],
    deps = [
        "//adhan_service:adhan_service",
        "//event_service:event_service",
        "//masjid_service:masjid_service",
        "//proto:adhan_service_go_proto",
        "//proto:event_service_go_proto",
        "//proto:masjid_service_go_proto",
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
    name = "test_data",
    testonly = 1,
    srcs = ["test_data.go"],
    importpath = "github.com/mnadev/limestone/test_infra",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:adhan_service_go_proto",
        "//proto:event_service_go_proto",
        "//proto:masjid_service_go_proto",
        "//proto:user_service_go_proto",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "unit_test_suite",
    testonly = 1,
    srcs = ["unit_test_suite.go"],
    importpath = "github.com/mnadev/limestone/test_infra",
    visibility = ["//visibility:public"],
    deps = [
        "//storage:storage_manager",
        "@com_github_proullon_ramsql//driver",
        "@com_github_stretchr_testify//suite",
        "@io_gorm_driver_postgres//:postgres",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//test/bufconn",
    ],
)

# Unit tests
go_test(
    name = "storage_manager_test",
    srcs = ["storage_manager_test.go"],
    embed = [
        ":assert",
        ":unit_test_suite",
        ":test_data",
    ],
    deps = [
        "//proto:adhan_service_go_proto",
        "//proto:event_service_go_proto",
        "//proto:masjid_service_go_proto",
        "//proto:user_service_go_proto",
        "//storage:storage_manager",
        "@com_github_stretchr_testify//suite",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)

# Integration tests
go_test(
    name = "adhan_service_test",
    srcs = ["adhan_service_test.go"],
    embed = [
        ":assert",
        ":integration_test_suite",
        ":test_data",
    ],
    deps = [
        "//proto:adhan_service_go_proto",
        "//adhan_service:adhan_service",
        "@com_github_stretchr_testify//suite",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)

go_test(
    name = "event_service_test",
    srcs = ["event_service_test.go"],
    embed = [
        ":assert",
        ":integration_test_suite",
        ":test_data",
    ],
    deps = [
        "//proto:event_service_go_proto",
        "//event_service:event_service",
        "@com_github_stretchr_testify//suite",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)

go_test(
    name = "masjid_service_test",
    srcs = ["masjid_service_test.go"],
    embed = [
        ":assert",
        ":integration_test_suite",
        ":test_data",
    ],
    deps = [
        "//proto:masjid_service_go_proto",
        "//masjid_service:masjid_service",
        "@com_github_stretchr_testify//suite",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)

go_test(
    name = "user_service_test",
    srcs = ["user_service_test.go"],
    embed = [
        ":assert",
        ":integration_test_suite",
        ":test_data",
    ],
    deps = [
        "//proto:user_service_go_proto",
        "//user_service:user_service",
        "@com_github_stretchr_testify//suite",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
