plugins {
	java
	war
	id("org.springframework.boot") version "3.0.4"
	//id("io.spring.dependency-management") version "1.1.0"
	id ("io.spring.dependency-management") version "1.0.11.RELEASE"
}

group = "com.javademoapp"
version = "0.0.1-SNAPSHOT"
java.sourceCompatibility = JavaVersion.VERSION_17

repositories {
	mavenCentral()
}

dependencies {
	implementation("org.springframework.boot:spring-boot-starter-web")
    testImplementation("junit:junit:4.13.1")
	implementation ("org.springframework.boot:spring-boot-starter-validation")
	developmentOnly("org.springframework.boot:spring-boot-devtools")
	providedRuntime("org.springframework.boot:spring-boot-starter-tomcat")
	testImplementation("org.springframework.boot:spring-boot-starter-test")
	implementation("com.fluxninja.aperture:aperture-java-core:0.26.0")

	configurations {
		all {
			exclude(group = "org.slf4j", module = "slf4j-simple")
		}
	}
	implementation("org.springframework.boot:spring-boot-starter-logging") {
		exclude(group = "org.slf4j", module = "slf4j-api")
	}
}

tasks.withType<Test> {
	useJUnitPlatform()
}
