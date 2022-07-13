# Bazel build file for SPS
load("@com_github_bazelbuild_buildtools//buildifier:def.bzl", "buildifier")

buildifier(
    name = "buildifier",
)

load("@io_bazel_rules_go//go:def.bzl", "go_binary")

#go_binary(
#    name = "test_bin",
#    srcs = ["test.go"],
#    importpath = "test",
#    visibility = ["//visibility:private"],
#)

go_binary(
    name = "sps",
    srcs = ["sps.go"],
    importpath = "test",
    visibility = ["//visibility:private"],
)
