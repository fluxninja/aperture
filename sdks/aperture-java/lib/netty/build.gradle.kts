plugins {
    id("aperture-java.java-library-conventions")
    id("aperture-java.publish-conventions")
}

dependencies {
    api(project(":lib:core"))
    api("io.netty:netty-all:4.1.41.Final")

    implementation("io.opentelemetry:opentelemetry-api:1.18.0")
}
