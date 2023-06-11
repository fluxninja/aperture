package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.*;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.google.protobuf.util.JsonFormat;
import io.opentelemetry.api.trace.Span;
import org.apache.http.HttpStatus;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** A Flow that can be accepted or rejected by Aperture Agent based on provided labels. */
public final class Flow {
    private final CheckResponse checkResponse;
    private final Span span;
    private boolean ended;
    private boolean failOpen;
    private FlowStatus flowStatus;

    private static final Logger logger = LoggerFactory.getLogger(Flow.class);

    Flow(CheckResponse checkResponse, Span span, boolean ended) {
        this.checkResponse = checkResponse;
        this.span = span;
        this.ended = ended;
        this.failOpen = true;
        this.flowStatus = FlowStatus.OK;
    }

    /**
     * Returns 'true' if flow was accepted by Aperture Agent, or if the Agent did not respond.
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
        if (this.checkResponse.getDecisionType()
                == CheckResponse.DecisionType.DECISION_TYPE_ACCEPTED) {
            return FlowDecision.Accepted;
        }
        return FlowDecision.Rejected;
    }

    /**
     * Returns whether the flow should be allowed to run, based on flow fail-open configuration and
     * Aperture Agent response. By default, flow will be allowed to run if Aperture Agent is
     * unreachable. To change this behavior, use {@link #withNoFailOpen()}.
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
     * @return This Flow object
     */
    public Flow withNoFailOpen() {
        this.failOpen = false;
        return this;
    }

    /**
     * Returns raw CheckResponse returned by Aperture Agent.
     *
     * @return raw CheckResponse returned by Aperture Agent.
     */
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

    /**
     * Set status of the flow to be ended. Primarily used in case of business logic failure after
     * the flow was accepted by Aperture Agent.
     *
     * @param status Status of the flow to be finished.
     */
    public void setStatus(FlowStatus status) {
        if (this.ended) {
            logger.warn("Trying to change status of an already ended flow");
        }
        this.flowStatus = status;
    }

    /**
     * Ends the flow, notifying the Aperture Agent whether it succeeded. Flow's Status is assumed to
     * be "OK" and can be set using {@link #setStatus}.
     */
    public void end() {
        if (this.ended) {
            logger.warn("Trying to end an already ended flow with status " + this.flowStatus);
            return;
        }
        this.ended = true;

        String checkResponseJSONBytes = "";
        try {
            if (this.checkResponse != null) {
                checkResponseJSONBytes = JsonFormat.printer().print(this.checkResponse);
            }
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
            logger.warn("Could not attach check response when ending flow", e);
        }

        logger.debug("Ending flow with status " + this.flowStatus);

        this.span
                .setAttribute(FLOW_STATUS_LABEL, this.flowStatus.name())
                .setAttribute(CHECK_RESPONSE_LABEL, checkResponseJSONBytes)
                .setAttribute(FLOW_STOP_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        this.span.end();
    }
}
