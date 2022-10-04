package com.fluxninja

import com.fluxninja.plugins.configureRouting
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.server.testing.*
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.time.Duration.Companion.milliseconds

class ApplicationTest {
    @Test
    fun testRoot() = testApplication {
        application {
            configureRouting(10, 50.milliseconds, 0.05F)
        }
        val response = client.get("/request") {}
        assertEquals(HttpStatusCode.OK, response.status)
    }
}
