plugins {
    id("java")
    id("application")
    id("com.github.johnrengelman.shadow")
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
    archiveBaseName.set("fatApp")
    archiveClassifier.set("")

    relocate("javassist", "com.exmaple.javassist")
}

dependencies {
    implementation("com.sparkjava:spark-core:2.9.4")
    implementation("io.grpc:grpc-api:1.44.0")
    implementation(project(":lib:core"))
}
