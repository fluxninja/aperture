plugins {
    id("java")
    id("application")
    id("com.github.johnrengelman.shadow")
}

application {
    mainClass.set("com.fluxninja.example.NettyServer")
}

dependencies {
    implementation(project(":lib:core"))
    implementation(project(":lib:netty"))
}
