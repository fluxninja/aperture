plugins {
    id("java")
    id("application")
    id("com.github.johnrengelman.shadow")
}

java {
    setSourceCompatibility("1.8")
    setTargetCompatibility("1.8")
}

application {
    mainClass.set("com.fluxninja.example.App")
}

tasks.jar {
    manifest {
        attributes["Main-Class"] = "com.fluxninja.example.App"
    }
}

tasks.shadowJar {
    archiveBaseName.set("standalone-example")
    archiveClassifier.set("all")

    mergeServiceFiles()
    relocate("javassist", "com.example.javassist")
}

dependencies {
    implementation("com.sparkjava:spark-core:2.9.4")
    implementation("io.grpc:grpc-api:1.44.0")
    implementation(project(":lib:core"))
}
