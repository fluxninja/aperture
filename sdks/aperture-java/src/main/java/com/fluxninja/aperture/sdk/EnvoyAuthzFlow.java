package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.envoy.service.auth.v3.CheckResponse;
import com.google.protobuf.util.JsonFormat;
import io.opentelemetry.api.trace.Span;

import static com.fluxninja.aperture.sdk.Constants.*;

public class EnvoyAuthzFlow {
    private final CheckResponse checkResponse;
    private final Span span;
    private boolean ended;

    EnvoyAuthzFlow(
            CheckResponse checkResponse,
            Span span,
            boolean ended) {
        this.checkResponse = checkResponse;
        this.span = span;
        this.ended = ended;
    }

    public boolean accepted() {
        if (this.checkResponse == null || !this.checkResponse.hasStatus()) {
            return true;
        }
        if (this.checkResponse.hasDeniedResponse()) {
            return false;
        }
        return this.checkResponse.getStatus().getCode() == 200;
    }

    public CheckResponse checkResponse() {
        return this.checkResponse;
    }

    public void end(FlowStatus statusCode) throws ApertureSDKException {
        if (this.ended) {
            throw new ApertureSDKException("Flow already ended");
        }
        this.ended = true;

        String checkResponseJSONBytes;
        try {
            checkResponseJSONBytes = JsonFormat.printer().print(this.checkResponse);
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
            throw new ApertureSDKException(e);
        }

        this.span.setAttribute(FEATURE_STATUS_LABEL, statusCode.name())
                .setAttribute(CHECK_RESPONSE_LABEL, checkResponseJSONBytes)
                .setAttribute(FLOW_STOP_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        this.span.end();
    }
}
