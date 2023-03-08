pluginManagement {
    plugins {
        id("com.google.protobuf") version "0.8.17"
        id("io.github.gradle-nexus.publish-plugin") version "1.1.0"
        id("com.github.johnrengelman.shadow") version "7.1.2"
    }
}
dependencyResolutionManagement {
    repositories {
        mavenCentral()
        mavenLocal()
        gradlePluginPortal()
    }
}

rootProject.name = "aperture-java"
include("lib:core", "lib:armeria", "lib:netty", "lib:servlet")
include("examples:armeria-example", "examples:spring-example", "examples:standalone-example", "examples:tomcat-example", "examples:netty-example")
include("javaagent:agent")
include("javaagent:test-services:armeria-test-service", "javaagent:test-services:netty-test-service")
