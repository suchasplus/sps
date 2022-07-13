# Bazel build file for SPS

load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "test_bin",
    srcs = ["test.go"],
    importpath = "test",
    visibility = ["//visibility:private"],
)
