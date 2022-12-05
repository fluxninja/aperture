plugins {
    id("aperture-java.java-library-conventions")
    id("aperture-java.publish-conventions")
}

dependencies {
    api("com.linecorp.armeria:armeria:1.15.0")
    api(project(":lib:core"))

    implementation("io.opentelemetry:opentelemetry-api:1.18.0")
}
