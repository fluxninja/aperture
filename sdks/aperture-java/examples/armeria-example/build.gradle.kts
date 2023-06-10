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
    implementation(project(":lib:core"))
    implementation(project(":lib:armeria"))

    runtimeOnly("org.slf4j:slf4j-simple:1.7.0")
}
