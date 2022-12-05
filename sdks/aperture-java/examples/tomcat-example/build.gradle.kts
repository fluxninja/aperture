plugins {
    id("java")
    id("war")
    id("com.bmuschko.tomcat") version "2.7.0"
}

repositories {
    mavenCentral()
    mavenLocal()
}

val tomcatVersion = "7.0.76"
dependencies {
    tomcat("org.apache.tomcat.embed:tomcat-embed-core:${tomcatVersion}")
    tomcat("org.apache.tomcat.embed:tomcat-embed-logging-juli:${tomcatVersion}")
    tomcat("org.apache.tomcat.embed:tomcat-embed-jasper:${tomcatVersion}")

    implementation(project(":lib:tomcat7"))

    providedCompile("javax.servlet:javax.servlet-api:3.1.0")
}

tomcat {
    contextPath = "/"
    httpPort = 58090
    httpsPort = 58091
}
