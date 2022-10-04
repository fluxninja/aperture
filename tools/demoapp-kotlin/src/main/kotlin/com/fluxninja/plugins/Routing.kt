package com.fluxninja.plugins

import io.ktor.server.routing.*
import io.ktor.server.application.*
import com.fluxninja.routes.*
import kotlin.time.Duration

fun Application.configureRouting(concurrency: Int, latency: Duration, rejectRation: Float) {
    routing {
        requestRoute(concurrency, latency, rejectRation)
    }
}
