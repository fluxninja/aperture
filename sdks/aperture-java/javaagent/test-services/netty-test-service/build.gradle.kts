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
    implementation("io.netty:netty-all:4.1.41.Final")
    implementation("org.slf4j:slf4j-simple:1.7.32")
}
