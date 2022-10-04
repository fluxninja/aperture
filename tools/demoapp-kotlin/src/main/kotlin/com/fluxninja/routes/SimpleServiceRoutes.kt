package com.fluxninja.routes

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.request.*
import com.fluxninja.models.*
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlin.random.Random
import kotlin.time.Duration
import kotlin.time.DurationUnit
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.plugins.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import kotlinx.coroutines.*

fun Route.requestRoute(concurrency: Int, latency: Duration, rejectRation: Float) {
    get("/request") {
        if (rejectRation > 0 && Random.nextFloat() < rejectRation) {
            call.respond(HttpStatusCode.ServiceUnavailable)
        }

        val requestBody = call.receive<Request>()
        val hostname = call.request.local.host

        var code = HttpStatusCode.OK
        for(chain in requestBody.chains) {
            if (chain.subrequests.isEmpty()) {
                call.respond(HttpStatusCode.BadRequest,"Received empty subrequest chain")
            }

            val requestDestination = chain.subrequests[0].destination
            if (requestDestination != hostname) {
                call.respond(HttpStatusCode.BadRequest,"Invalid message destination")
            }

            code = processChain(chain, concurrency, latency)
        }

        call.respond(code)
    }
}

fun processChain(chain: SubrequestChain, concurrency: Int, latency: Duration): HttpStatusCode {
    if (chain.subrequests.size == 1) {
        return processRequest(chain.subrequests[0], concurrency, latency)
    }
    val requestForwardingDestination = chain.subrequests[1].destination
    val trimmedSubrequestChain = SubrequestChain(subrequests = chain.subrequests.slice(1 until chain.subrequests.size))
    val trimmedRequest = Request(chains = listOf(trimmedSubrequestChain))

    return forwardRequest(requestForwardingDestination, trimmedRequest)
}

fun forwardRequest(destinationHostname: String, requestBody: Request): HttpStatusCode {
    val jsonBody = Json.encodeToString(requestBody)

    runBlocking {
        val client = HttpClient(CIO) {
            defaultRequest {
                url {
                    protocol = URLProtocol.HTTPS
                    host = "ktor.io"
                    path("docs/")
                    parameters.append("token", "abc123")
                }
                header("X-Custom-Header", "Hello")
            }
        }
        val response: HttpResponse = client.get("welcome.html")
}

fun processRequest(subrequest: Subrequest, concurrency: Int, latency: Duration): HttpStatusCode {
    if (concurrency > 0) {

    }

    l = latency.toLong(DurationUnit.MILLISECONDS)
    if (l > 0) {
        Thread.sleep(l)
    }

    return HttpStatusCode.OK
}
