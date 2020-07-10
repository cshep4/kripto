load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

RULES_JVM_EXTERNAL_TAG = "2.8"
RULES_JVM_EXTERNAL_SHA = "79c9850690d7614ecdb72d68394f994fef7534b292c4867ce5e7dec0aa7bdfad"

http_archive(
    name = "rules_jvm_external",
    strip_prefix = "rules_jvm_external-%s" % RULES_JVM_EXTERNAL_TAG,
    sha256 = RULES_JVM_EXTERNAL_SHA,
    url = "https://github.com/bazelbuild/rules_jvm_external/archive/%s.zip" % RULES_JVM_EXTERNAL_TAG,
)

load("@rules_jvm_external//:defs.bzl", "maven_install")

maven_install(
    name = "maven",
    artifacts = [
        "org.jetbrains.kotlin:kotlin-stdlib-jdk8:1.3.72",
        "org.litote.kmongo:kmongo:4.0.1",
    ],
    repositories = [
        "https://jcenter.bintray.com/",
    ],
    fetch_sources = True,   # Fetch source jars. Defaults to False.
)

maven_install(
    name = "re_deps",
    artifacts = [
        "org.jetbrains.kotlin:kotlin-stdlib-jdk8:1.3.72",
        "org.litote.kmongo:kmongo:4.0.1",
        "com.amazonaws:aws-java-sdk-lambda:1.11.819",
        "com.amazonaws:aws-lambda-java-events:3.1.0",
        "com.amazonaws:aws-lambda-java-core:1.2.1",
        "com.google.code.gson:gson:2.8.6",
        "io.lumigo:java-tracer:1.0.30",
        "io.lumigo:lumigo-agent:1.0.30",
        "com.sendgrid:sendgrid-java:4.5.0",
    ],
    repositories = [
        "https://maven.google.com",
        "https://jcenter.bintray.com/",
        "http://maven.nuiton.org/release",
    ],
    fetch_sources = True,   # Fetch source jars. Defaults to False.
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

rules_kotlin_version = "legacy-1.3.0"
rules_kotlin_sha = "4fd769fb0db5d3c6240df8a9500515775101964eebdf85a3f9f0511130885fde"
http_archive(
    name = "io_bazel_rules_kotlin",
    urls = ["https://github.com/bazelbuild/rules_kotlin/archive/%s.zip" % rules_kotlin_version],
    type = "zip",
    strip_prefix = "rules_kotlin-%s" % rules_kotlin_version,
    sha256 = rules_kotlin_sha,
)

load("@io_bazel_rules_kotlin//kotlin:kotlin.bzl", "kotlin_repositories", "kt_register_toolchains")
kotlin_repositories() # if you want the default. Otherwise see custom kotlinc distribution below
kt_register_toolchains() # to use the default toolchain, otherwise see toolchains below