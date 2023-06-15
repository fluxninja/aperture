plugins {
	java
	war
	id("org.springframework.boot") version "2.7.9"
	id("io.spring.dependency-management") version "1.0.11.RELEASE"
}

group = "com.javademoapp"
version = "0.0.1-SNAPSHOT"

repositories {
	mavenCentral()
}

dependencies {
	configurations {
        all {
            exclude(group = "org.slf4j", module = "slf4j-simple")
        }
    }
    implementation("org.springframework.boot:spring-boot-starter-logging") {
            exclude(group = "org.slf4j", module = "slf4j-api")
    }

	implementation("org.springframework.boot:spring-boot-starter-actuator:2.7.9")
	implementation("io.micrometer:micrometer-registry-prometheus:1.9.0")
	implementation("org.springframework.boot:spring-boot-starter-web:2.7.9")
	implementation("com.fluxninja.aperture:aperture-java-core:2.3.0")
	implementation("com.fluxninja.aperture:aperture-java-servlet:2.3.0")

	developmentOnly("org.springframework.boot:spring-boot-devtools:2.7.9")
	providedRuntime("org.springframework.boot:spring-boot-starter-tomcat:2.7.9")

	testImplementation("junit:junit:4.13.1")
	testImplementation("org.springframework.boot:spring-boot-starter-test")
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}
