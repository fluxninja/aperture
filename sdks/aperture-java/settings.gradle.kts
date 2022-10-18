pluginManagement {
    plugins {
        id("com.google.protobuf") version "0.8.17"
        id("io.github.gradle-nexus.publish-plugin") version "1.1.0"
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
