package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.*;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.google.protobuf.util.JsonFormat;
import io.opentelemetry.api.trace.Span;
import org.apache.http.HttpStatus;

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
     * @deprecated This method assumes fail-open behavior. Use {@link #getDecision} instead
     * @return Whether the flow was accepted.
     */
    public boolean accepted() {
        return getDecision() == FlowDecision.Unreachable || getDecision() == FlowDecision.Accepted;
    }

    /**
     * Returns Aperture Agent's decision or information on Agent being unreachable.
     *
     * @return Result of Check query
     */
    public FlowDecision getDecision() {
        if (this.checkResponse == null) {
            return FlowDecision.Unreachable;
        }
        if (this.checkResponse.getDecisionType()
                == CheckResponse.DecisionType.DECISION_TYPE_ACCEPTED) {
            return FlowDecision.Accepted;
        }
        return FlowDecision.Rejected;
    }

    public CheckResponse checkResponse() {
        return this.checkResponse;
    }

    /**
     * Returns Aperture Agent's reason for rejecting the flow. Reason is represented by an
     * appropriate HTTP code. If the flow was not rejected, an IllegalStateException will be thrown.
     *
     * @return HTTP code of rejection reason
     */
    public int getRejectionHttpStatusCode() {
        if (this.getDecision() == FlowDecision.Rejected) {
            switch (this.checkResponse.getRejectReason()) {
                case REJECT_REASON_RATE_LIMITED:
                    return HttpStatus.SC_TOO_MANY_REQUESTS;
                case REJECT_REASON_NO_TOKENS:
                    return HttpStatus.SC_SERVICE_UNAVAILABLE;
                case REJECT_REASON_REGULATED:
                    return HttpStatus.SC_FORBIDDEN;
                default:
                    return HttpStatus.SC_FORBIDDEN;
            }
        } else {
            throw new IllegalStateException("Flow not rejected");
        }
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
