package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.FlowEndResponse;

public class EndResponse {
    private final FlowEndResponse flowEndResponse;
    private final Exception error;

    public EndResponse(FlowEndResponse flowEndResponse, Exception error) {
        this.flowEndResponse = flowEndResponse;
        this.error = error;
    }

    public FlowEndResponse getFlowEndResponse() {
        return flowEndResponse;
    }

    public Exception getError() {
        return error;
    }
}
