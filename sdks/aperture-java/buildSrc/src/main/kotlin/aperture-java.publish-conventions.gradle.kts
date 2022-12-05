plugins {
    `maven-publish`
    signing
}

repositories {
    mavenCentral()
    gradlePluginPortal()
    mavenLocal()
}

publishing {
    publications {
        register<MavenPublication>("mavenPublication") {
            groupId = "com.fluxninja.aperture"
            artifactId = "aperture-java-" + project.name

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
