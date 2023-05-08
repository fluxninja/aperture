package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.*;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.google.protobuf.util.JsonFormat;
import io.opentelemetry.api.trace.Span;

public final class Flow {
    private final CheckResponse checkResponse;
    private final Span span;
    private boolean ended;

    Flow(CheckResponse checkResponse, Span span, boolean ended) {
        this.checkResponse = checkResponse;
        this.span = span;
        this.ended = ended;
    }

    /**
     * Returns 'true' if flow was accepted by Aperture Agent, or if the Agent did not respond.
     *
     * @deprecated This method assumes fail-open behavior. Use {@link #result} instead
     * @return Whether the flow was accepted.
     */
    public boolean accepted() {
        return result() == FlowResult.Unreachable || result() == FlowResult.Accepted;
    }

    /**
     * Returns Aperture Agent's decision or information on Agent being unreachable.
     *
     * @return Result of Check query
     */
    public FlowResult result() {
        if (this.checkResponse == null) {
            return FlowResult.Unreachable;
        }
        if (this.checkResponse.getDecisionType()
                == CheckResponse.DecisionType.DECISION_TYPE_ACCEPTED) {
            return FlowResult.Accepted;
        }
        return FlowResult.Rejected;
    }

    public CheckResponse checkResponse() {
        return this.checkResponse;
    }

    public void end(FlowStatus statusCode) throws ApertureSDKException {
        if (this.ended) {
            throw new ApertureSDKException("Flow already ended");
        }
        this.ended = true;

        String checkResponseJSONBytes = "";
        try {
            if (this.checkResponse != null) {
                checkResponseJSONBytes = JsonFormat.printer().print(this.checkResponse);
            }
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
            throw new ApertureSDKException(e);
        }

        this.span
                .setAttribute(FLOW_STATUS_LABEL, statusCode.name())
                .setAttribute(CHECK_RESPONSE_LABEL, checkResponseJSONBytes)
                .setAttribute(FLOW_STOP_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        this.span.end();
    }
}
