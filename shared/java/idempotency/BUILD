package(default_visibility = ["//visibility:public"])
load("@io_bazel_rules_kotlin//kotlin:kotlin.bzl", "kt_jvm_library")

kt_jvm_library(
    name = "java_idempotency",
    srcs = glob(["src/main/kotlin/com/cshep4/kripto/idempotency/**"]),
    deps = [
            "@maven//:org_jetbrains_kotlin_kotlin_stdlib_jdk8",
            "@maven//:org_litote_kmongo_kmongo",
        ],
)