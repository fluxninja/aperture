package com.fluxninja.routes

import com.fluxninja.models.Request
import com.fluxninja.models.SubrequestChain
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.coroutines.runBlocking
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlin.random.Random
import kotlin.time.Duration
import kotlin.time.DurationUnit

fun Route.requestRoute(concurrency: Int, latency: Duration, rejectRation: Float) {
    post("/request") {
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
        return processRequest(concurrency, latency)
    }
    val requestForwardingDestination = chain.subrequests[1].destination
    val trimmedSubrequestChain = SubrequestChain(subrequests = chain.subrequests.slice(1 until chain.subrequests.size-1))
    val trimmedRequest = Request(chains = listOf(trimmedSubrequestChain))

    return forwardRequest(requestForwardingDestination, trimmedRequest)
}

fun forwardRequest(destinationHostname: String, requestBody: Request): HttpStatusCode {
    val jsonBody = Json.encodeToString(requestBody)

    var code = HttpStatusCode.OK
    runBlocking {
        val client = HttpClient(CIO)
        val response: HttpResponse = client.post(destinationHostname) {
            contentType(ContentType.Application.Json)
            setBody(jsonBody)
        }
        code = response.status
    }

    return code
}

fun processRequest(concurrency: Int, latency: Duration): HttpStatusCode {
    if (concurrency > 0) {

    }

    val l = latency.toLong(DurationUnit.MILLISECONDS)
    if (l > 0) {
        Thread.sleep(l)
    }

    return HttpStatusCode.OK
}
