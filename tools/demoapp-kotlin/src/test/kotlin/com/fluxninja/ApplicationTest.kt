package com.fluxninja

import com.fluxninja.plugins.configureRouting
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.server.testing.*
import kotlin.test.Test
import kotlin.test.assertEquals

class ApplicationTest {
    @Test
    fun testRoot() = testApplication {
        application {
            configureRouting()
        }
        client.get("/request").apply {
            assertEquals(HttpStatusCode.OK, status)
        }
    }
}
