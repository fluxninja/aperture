plugins {
    id("aperture-java.java-library-conventions")
    id("aperture-java.publish-conventions")
}

dependencies {
    api(project(":lib:core"))
    api("jakarta.servlet:jakarta.servlet-api:6.0.0")
    api("javax.servlet:javax.servlet-api:3.1.0")

    implementation("io.opentelemetry:opentelemetry-api:1.18.0")
}
