plugins {
    id("java")
    id("application")
}

application {
    mainClass.set("com.fluxninja.example.App")
}

dependencies {
    implementation("com.sparkjava:spark-core:2.9.4")
    implementation("io.grpc:grpc-api:1.44.0")
    implementation(project(":lib:core"))
}
