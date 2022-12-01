import com.github.jengelman.gradle.plugins.shadow.tasks.ConfigureShadowRelocation

plugins {
    id("com.github.johnrengelman.shadow")
    id("java-library")
    id("com.google.protobuf")

    `maven-publish`
    signing
}

val relocateShadowJar = tasks.register<ConfigureShadowRelocation>("relocateShadowJar") {
    target = tasks.shadowJar.get()
    prefix = "apertureshadow"
}
tasks.shadowJar.get().dependsOn(relocateShadowJar.get())
tasks.shadowJar {
    mergeServiceFiles()
}

subprojects {
    group = "com.fluxninja.aperture"
}

dependencies {
    api("com.linecorp.armeria:armeria:1.15.0")

    implementation(platform("io.opentelemetry:opentelemetry-bom-alpha:1.18.0-alpha"))
    implementation("io.opentelemetry:opentelemetry-sdk-trace:1.18.0")
    implementation("io.opentelemetry:opentelemetry-exporter-otlp:1.18.0")
    implementation("io.opentelemetry:opentelemetry-exporter-logging:1.18.0")
    implementation("io.grpc:grpc-protobuf:1.44.0")
    implementation("io.grpc:grpc-stub:1.44.0")
    implementation("org.slf4j:slf4j-simple:1.7.0")
    implementation("com.google.protobuf:protobuf-java-util:3.21.6")

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

publishing {
    publications {
        register<MavenPublication>("mavenPublication") {
            groupId = "com.fluxninja.aperture"
            artifactId = "aperture-java"

            from(components["java"])

            versionMapping {
                allVariants {
                    fromResolutionResult()
                }
            }

            pom {
                name.set("Aperture Java")
                description.set("Java SDK to connect to FluxNinja Aperture")
                url.set("https://github.com/fluxninja/aperture-java")

                licenses {
                    license {
                        name.set("The Apache License, Version 2.0")
                        url.set("http://www.apache.org/licenses/LICENSE-2.0.txt")
                    }
                }

                developers {
                    developer {
                        id.set("fluxninja")
                        name.set("FluxNinja")
                        url.set("https://github.com/fluxninja")
                    }
                }

                scm {
                    connection.set("scm:git:git@github.com:fluxninja/aperture-java.git")
                    developerConnection.set("scm:git:git@github.com:fluxninja/aperture-java.git")
                    url.set("git@github.com:fluxninja/aperture-java.git")
                }
            }
        }
    }
}

if (System.getenv("CI") != null) {
    signing {
        useInMemoryPgpKeys(System.getenv("GPG_PRIVATE_KEY"), System.getenv("GPG_PASSWORD"))
        sign(publishing.publications["mavenPublication"])
    }
}
