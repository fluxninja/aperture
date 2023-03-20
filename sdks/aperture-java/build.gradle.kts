import java.time.Duration

plugins {
    id("io.github.gradle-nexus.publish-plugin")
    id("com.diffplug.gradle.spotless")
}

subprojects {
    group = "com.fluxninja.aperture"
}

apply(from = "version.gradle.kts")

nexusPublishing {
    packageGroup.set("com.fluxninja.aperture")

    repositories {
        sonatype {
            nexusUrl.set(uri("https://s01.oss.sonatype.org/service/local/"))
            snapshotRepositoryUrl.set(uri("https://s01.oss.sonatype.org/content/repositories/snapshots/"))
            username.set(System.getenv("SONATYPE_USER"))
            password.set(System.getenv("SONATYPE_KEY"))
        }
    }

    connectTimeout.set(Duration.ofMinutes(5))
    clientTimeout.set(Duration.ofMinutes(5))

    transitionCheckOptions {
        // We have many artifacts so Maven Central takes a long time on its compliance checks. This sets
        // the timeout for waiting for the repository to close to a comfortable 50 minutes.
        maxRetries.set(300)
        delayBetween.set(Duration.ofSeconds(10))
    }
}

spotless {
    java {
        target("**/*.java")
        targetExclude("lib/core/src/main/java/com/fluxninja/generated/**/*.java")

        googleJavaFormat().aosp()
    }
}
