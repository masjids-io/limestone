load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "adhan_file",
    srcs = ["adhan_file.go"],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:adhan_service_go_proto",
        "@com_github_google_uuid//:uuid",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "events",
    srcs = ["events.go"],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:event_service_go_proto",
        "@com_github_google_uuid//:uuid",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "masjid",
    srcs = [
        "masjid.go",
        "prayer_config.go",
    ],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:masjid_service_go_proto",
        "@com_github_google_uuid//:uuid",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "prayer_config",
    srcs = ["prayer_config.go"],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//proto:masjid_service_go_proto",
    ],
)

go_library(
    name = "storage_manager",
    srcs = [
        "adhan_file.go",
        "events.go",
        "masjid.go",
        "prayer_config.go",
        "storage_manager.go",
        "user.go",
    ],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//auth:password",
        "//proto:adhan_service_go_proto",
        "//proto:event_service_go_proto",
        "//proto:masjid_service_go_proto",
        "//proto:user_service_go_proto",
        "@com_github_google_uuid//:uuid",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "user",
    srcs = ["user.go"],
    importpath = "github.com/mnadev/limestone/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//auth:password",
        "//proto:user_service_go_proto",
        "@com_github_google_uuid//:uuid",
        "@io_gorm_gorm//:gorm",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
