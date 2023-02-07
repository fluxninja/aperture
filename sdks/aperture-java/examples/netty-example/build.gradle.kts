plugins {
    id("java")
    id("application")
    id("com.github.johnrengelman.shadow")
}

tasks.shadowJar {
    mergeServiceFiles()
}

java {
    setSourceCompatibility("1.8")
    setTargetCompatibility("1.8")
}

application {
    mainClass.set("com.fluxninja.example.NettyServer")
}

dependencies {
    implementation(project(":lib:core"))
    implementation(project(":lib:netty"))
}
