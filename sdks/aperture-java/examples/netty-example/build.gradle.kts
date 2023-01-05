plugins {
    id("java")
    id("application")
}

application {
    mainClass.set("com.fluxninja.example.NettyServer")
}

dependencies {
    implementation(project(":lib:core"))
    implementation(project(":lib:netty"))
}
