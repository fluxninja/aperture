plugins {
    id("java-library")
}

java {
    setSourceCompatibility("1.8")
    setTargetCompatibility("1.8")
    withJavadocJar()
    withSourcesJar()
}

repositories {
    mavenCentral()
    gradlePluginPortal()
    mavenLocal()
}
