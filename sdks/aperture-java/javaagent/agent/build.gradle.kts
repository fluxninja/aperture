plugins {
    `maven-publish`
    signing
    id("java")
    id("com.github.johnrengelman.shadow")
    id("org.springframework.boot") version "2.7.9"
    id("io.spring.dependency-management") version "1.1.0"
}

java {
    setSourceCompatibility("1.8")
    setTargetCompatibility("1.8")
}

tasks.assemble.get().dependsOn(tasks.shadowJar.get())

tasks.jar {
    manifest {
        attributes["Premain-Class"] = "com.fluxninja.aperture.instrumentation.ApertureInstrumentationAgent"
    }
    enabled = false
}

tasks.shadowJar {
    archiveBaseName.set("agent")
    archiveClassifier.set("")

    mergeServiceFiles()
    relocate("javassist", "com.example.javassist")
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
            artifactId = "aperture-javaagent"

            artifact(tasks["shadowJar"])
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


dependencies {
    implementation("net.bytebuddy:byte-buddy-dep:1.12.19")
    implementation("org.slf4j:log4j-over-slf4j:2.0.6")
    implementation("org.springframework.boot:spring-boot-starter-web") {
        exclude("org.springframework.boot", "spring-boot-starter-logging")
    }

    implementation(project(":lib:core"))
    implementation(project(":lib:armeria"))
    implementation(project(":lib:netty"))
}

repositories {
    mavenCentral()
}
