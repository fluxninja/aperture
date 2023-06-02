package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.*;

import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse;
import com.fluxninja.generated.google.rpc.Status;
import com.google.protobuf.Value;
import com.google.protobuf.util.JsonFormat;
import com.google.rpc.Code;
import io.opentelemetry.api.trace.Span;
import org.apache.http.HttpStatus;

public class TrafficFlow {
    private final CheckHTTPResponse checkResponse;
    private final Span span;
    public boolean ended;
    private boolean ignored;
    private boolean failOpen;

    TrafficFlow(CheckHTTPResponse checkResponse, Span span, boolean ended) {
        this.checkResponse = checkResponse;
        this.span = span;
        this.ended = ended;
        this.ignored = false;
        this.failOpen = true;
    }

    static TrafficFlow ignoredFlow() {
        TrafficFlow flow = new TrafficFlow(successfulResponse(), null, true);
        flow.ignored = true;
        return flow;
    }

    /**
     * Returns 'true' if flow was accepted by Aperture Agent, or if the agent did not respond.
     *
     * @deprecated This method assumes fail-open behavior. Use {@link #shouldRun} or {@link
     *     #getDecision} instead
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
        if (this.checkResponse.getStatus().getCode() == Code.OK_VALUE) {
            return FlowDecision.Accepted;
        }
        return FlowDecision.Rejected;
    }

    /**
     * Returns whether the flow should be allowed to run, based on flow fail-open configuration and
     * Aperture Agent response.
     *
     * @return Whether the flow should be allowed to run
     */
    public boolean shouldRun() {
        return getDecision() == FlowDecision.Accepted
                || (getDecision() == FlowDecision.Unreachable && this.failOpen);
    }

    /**
     * Disables fail-open behavior. If set, the {@link #shouldRun} method will return False if the
     * Aperture Agent is unreachable.
     *
     * @return This TrafficFlow object
     */
    public TrafficFlow withNoFailOpen() {
        this.failOpen = false;
        return this;
    }

    public boolean ignored() {
        return this.ignored;
    }

    public CheckHTTPResponse checkResponse() {
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
            if (this.checkResponse().hasDeniedResponse()
                    && this.checkResponse().getDeniedResponse().getStatus() != 0) {
                return this.checkResponse.getDeniedResponse().getStatus();
            } else {
                return HttpStatus.SC_FORBIDDEN;
            }
        } else {
            throw new IllegalStateException("Flow not rejected");
        }
    }

    public void end(FlowStatus statusCode) throws ApertureSDKException {
        if (this.ignored) {
            // span has not been started, and so doesn't need to be ended.
            return;
        }
        if (this.ended) {
            throw new ApertureSDKException("Flow already ended");
        }
        this.ended = true;

        String serializedFlowcontrolCheckResponse = "";
        if (this.checkResponse != null
                && this.checkResponse.hasDynamicMetadata()
                && this.checkResponse
                        .getDynamicMetadata()
                        .getFieldsMap()
                        .containsKey("aperture.check_response")) {
            Value checkResponse =
                    this.checkResponse
                            .getDynamicMetadata()
                            .getFieldsMap()
                            .get("aperture.check_response");
            if (checkResponse.hasStringValue()) {
                // If checkResponse comes pre-serialized from envoy, pass it
                // through as-is.
                serializedFlowcontrolCheckResponse = checkResponse.getStringValue();
            } else {
                // Otherwise, serialize it.
                try {
                    serializedFlowcontrolCheckResponse = JsonFormat.printer().print(checkResponse);
                } catch (com.google.protobuf.InvalidProtocolBufferException e) {
                    throw new ApertureSDKException(e);
                }
            }
        }

        this.span
                .setAttribute(FLOW_STATUS_LABEL, statusCode.name())
                .setAttribute(CHECK_RESPONSE_LABEL, serializedFlowcontrolCheckResponse)
                .setAttribute(FLOW_STOP_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        this.span.end();
    }

    // Artificial response if none is received from agent
    static CheckHTTPResponse successfulResponse() {
        return CheckHTTPResponse.newBuilder()
                .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
                .build();
    }
}
