package com.fluxninja

import io.ktor.server.engine.*
import io.ktor.server.netty.*
import com.fluxninja.plugins.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.plugins.contentnegotiation.*
import kotlinx.serialization.json.Json
import kotlin.time.Duration
import kotlin.time.Duration.Companion.milliseconds

fun main() {
    val concurrency: Int = (System.getenv("SIMPLE_SERVICE_CONCURRENCY") ?: 10) as Int
    val latency: Duration = (System.getenv("SIMPLE_SERVICE_LATENCY") ?: 50.milliseconds) as Duration
    val rejectRatio: Float = (System.getenv("SIMPLE_SERVICE_REJECT_RATIO") ?: 0.05F) as Float

    embeddedServer(Netty, port = 8080, host = "0.0.0.0") {
        install(ContentNegotiation) {
            json(Json {
                prettyPrint = true
            })
        }
        configureSerialization()
        configureRouting(concurrency, latency, rejectRatio)
    }.start(wait = true)
}
