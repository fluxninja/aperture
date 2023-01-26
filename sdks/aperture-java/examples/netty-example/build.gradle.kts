plugins {
    id("java")
    id("application")
    id("com.github.johnrengelman.shadow")
}

tasks.shadowJar {
    mergeServiceFiles()
}

application {
    mainClass.set("com.fluxninja.example.NettyServer")
}

dependencies {
    implementation(project(":lib:core"))
    implementation(project(":lib:netty"))
}
