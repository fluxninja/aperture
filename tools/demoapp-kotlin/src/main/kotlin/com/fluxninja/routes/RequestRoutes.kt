package com.fluxninja.routes

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.requestRoute() {
    get("/request") {
        call.respond(HttpStatusCode.OK)
    }
}
