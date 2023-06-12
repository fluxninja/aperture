plugins {
	id("java")
	id("application")
	id("com.github.johnrengelman.shadow")
	id("org.springframework.boot") version "2.7.9"
	id("io.spring.dependency-management") version "1.1.0"
}

tasks.shadowJar {
	mergeServiceFiles()
}

application {
	mainClass.set("com.fluxninja.example.SpringBootApp")
}

java.sourceCompatibility = JavaVersion.VERSION_1_8

repositories {
	mavenCentral()
}

dependencies {
	implementation("org.springframework.boot:spring-boot-starter-web") {
		exclude("org.springframework.boot", "spring-boot-starter-logging")
	}

	implementation(project(":lib:core"))
	implementation(project(":lib:servlet"))

	runtimeOnly("org.slf4j:slf4j-simple:1.7.0")
}
