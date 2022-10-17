import java.time.Duration

plugins {
    id("application")
    id("java")
    id("com.google.protobuf")
    id("io.github.gradle-nexus.publish-plugin")

    `maven-publish`
    signing
}

application {
    mainClass.set("com.fluxninja.aperture.example.App")
}

apply(from = "version.gradle.kts")

nexusPublishing {
    packageGroup.set("com.fluxninja.aperture")

    repositories {
        sonatype {
            nexusUrl.set(uri("https://s01.oss.sonatype.org/service/local/"))
            snapshotRepositoryUrl.set(uri("https://s01.oss.sonatype.org/content/repositories/snapshots/"))
            username.set(System.getenv("SONATYPE_USER"))
            password.set(System.getenv("SONATYPE_KEY"))
        }
    }

    connectTimeout.set(Duration.ofMinutes(5))
    clientTimeout.set(Duration.ofMinutes(5))

    transitionCheckOptions {
        // We have many artifacts so Maven Central takes a long time on its compliance checks. This sets
        // the timeout for waiting for the repository to close to a comfortable 50 minutes.
        maxRetries.set(300)
        delayBetween.set(Duration.ofSeconds(10))
    }
}

subprojects {
    group = "com.fluxninja.aperture"
}

dependencies {
    implementation(platform("io.opentelemetry:opentelemetry-bom-alpha:1.18.0-alpha"))
    implementation("io.opentelemetry:opentelemetry-sdk-trace:1.18.0")
    implementation("io.opentelemetry:opentelemetry-exporter-otlp-trace:1.14.0")
    implementation("io.opentelemetry:opentelemetry-exporter-logging:1.18.0")
    implementation("com.sparkjava:spark-core:2.9.4")
    implementation("io.grpc:grpc-protobuf:1.44.0")
    implementation("io.grpc:grpc-stub:1.44.0")
    implementation("org.slf4j:slf4j-simple:2.0.1")
    implementation("com.google.protobuf:protobuf-java-util:3.21.6")
    implementation("com.linecorp.armeria:armeria:1.20.0")

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
