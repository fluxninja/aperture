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
    mainClass.set("com.fluxninja.example.ArmeriaServer")
}

dependencies {
    implementation("com.linecorp.armeria:armeria:1.15.0")
}
