package com.fluxninja.models

import kotlinx.serialization.*

@Serializable
data class Request(var chains: List<SubrequestChain>)

@Serializable
data class SubrequestChain(var subrequests: List<Subrequest>)

@Serializable
data class Subrequest(var destination: String)
