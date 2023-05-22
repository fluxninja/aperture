import com.github.jengelman.gradle.plugins.shadow.tasks.ConfigureShadowRelocation

plugins {
    id("com.github.johnrengelman.shadow")
    id("aperture-java.java-library-conventions")
    id("aperture-java.publish-conventions")

    id("com.google.protobuf")
}

val relocateShadowJar = tasks.register<ConfigureShadowRelocation>("relocateShadowJar") {
    target = tasks.shadowJar.get()
    prefix = "apertureshadow"
}
tasks.shadowJar.get().dependsOn(relocateShadowJar.get())
tasks.shadowJar {
    mergeServiceFiles()
}

dependencies {
    api("com.google.protobuf:protobuf-java-util:3.22.2")

    implementation(platform("io.opentelemetry:opentelemetry-bom-alpha:1.18.0-alpha"))
    implementation("io.opentelemetry:opentelemetry-sdk-trace:1.18.0")
    implementation("io.opentelemetry:opentelemetry-exporter-otlp:1.18.0")
    implementation("io.opentelemetry:opentelemetry-exporter-logging:1.18.0")
    implementation("io.grpc:grpc-protobuf:1.44.0")
    implementation("io.grpc:grpc-stub:1.44.0")
    implementation("org.apache.httpcomponents:httpcore:4.4.16")
    implementation("org.slf4j:slf4j-simple:1.7.0")

    runtimeOnly("io.grpc:grpc-netty-shaded:1.49.0")

    // Workaround for @javax.annotation.Generated
    // see: https://github.com/grpc/grpc-java/issues/3633
    compileOnly("javax.annotation:javax.annotation-api:1.3.2")
    compileOnly("io.grpc:grpc-api:1.44.0")

    testImplementation("org.junit.jupiter:junit-jupiter-api:5.8.1")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:5.8.1")
}

// Publishing
java {
    setSourceCompatibility("1.8")
    setTargetCompatibility("1.8")
    withJavadocJar()
    withSourcesJar()
}
