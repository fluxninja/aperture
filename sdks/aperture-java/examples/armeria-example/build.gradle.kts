plugins {
    id("java")
    id("application")
    id("com.github.johnrengelman.shadow")
}

tasks.shadowJar {
    mergeServiceFiles()
}

application {
    mainClass.set("com.fluxninja.example.ArmeriaServer")
}

dependencies {
    implementation("com.linecorp.armeria:armeria:1.15.0")
    implementation(project(":lib:core"))
    implementation(project(":lib:armeria"))
}
