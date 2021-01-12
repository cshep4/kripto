load("@io_bazel_rules_kotlin//kotlin:kotlin.bzl", "kt_jvm_library")

kt_jvm_library(
    name = "receipt_emailer",
    srcs = glob(["src/main/kotlin/com/cshep4/kripto/receiptemailer/**"]),
    deps = [
            "//shared/java/idempotency:java_idempotency",
            "@re_deps//:org_jetbrains_kotlin_kotlin_stdlib_jdk8",
            "@re_deps//:org_litote_kmongo_kmongo",
            "@re_deps//:com_amazonaws_aws_java_sdk_lambda",
            "@re_deps//:com_amazonaws_aws_lambda_java_events",
            "@re_deps//:com_amazonaws_aws_lambda_java_core",
            "@re_deps//:com_google_code_gson_gson",
            "@re_deps//:com_sendgrid_sendgrid_java",
        ],
)