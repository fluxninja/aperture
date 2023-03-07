plugins {
	id("java")
	id("application")
	id("com.github.johnrengelman.shadow")
	id("org.springframework.boot") version "3.0.2"
	id("io.spring.dependency-management") version "1.1.0"
}

tasks.shadowJar {
	mergeServiceFiles()
}

application {
	mainClass.set("com.fluxninja.example.SpringBootApp")
}

java.sourceCompatibility = JavaVersion.VERSION_17

repositories {
	mavenCentral()
}

dependencies {
	implementation("org.springframework.boot:spring-boot-starter-web") {
		exclude("org.springframework.boot", "spring-boot-starter-logging")
	}

	implementation(project(":lib:core"))
	implementation(project(":lib:servlet"))
}
