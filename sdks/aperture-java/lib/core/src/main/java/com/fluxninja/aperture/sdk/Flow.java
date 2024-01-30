package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.CHECK_RESPONSE_LABEL;
import static com.fluxninja.aperture.sdk.Constants.FLOW_STATUS_LABEL;
import static com.fluxninja.aperture.sdk.Constants.FLOW_STOP_TIMESTAMP_LABEL;

import com.fluxninja.aperture.sdk.cache.KeyDeleteResponse;
import com.fluxninja.aperture.sdk.cache.KeyLookupResponse;
import com.fluxninja.aperture.sdk.cache.KeyUpsertResponse;
import com.fluxninja.aperture.sdk.cache.LookupStatus;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.*;
import com.google.protobuf.ByteString;
import com.google.protobuf.util.JsonFormat;
import io.opentelemetry.api.trace.Span;
import java.time.Duration;
import java.util.ArrayList;
import org.apache.http.HttpStatus;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** A Flow that can be accepted or rejected by Aperture Agent based on provided labels. */
public final class Flow {
    private final CheckResponse checkResponse;
    private final Span span;
    private boolean ended;
    private boolean rampMode;
    private FlowStatus flowStatus;
    private FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient;
    private String cacheKey;
    private String controlPoint;
    private Exception error;

    private static final Logger logger = LoggerFactory.getLogger(Flow.class);

    Flow(
            CheckResponse checkResponse,
            Span span,
            boolean ended,
            boolean rampMode,
            FlowControlServiceGrpc.FlowControlServiceBlockingStub fcs,
            String resultCacheKey,
            String controlPoint,
            Exception error) {
        this.checkResponse = checkResponse;
        this.span = span;
        this.ended = ended;
        this.rampMode = rampMode;
        this.flowStatus = FlowStatus.OK;
        this.flowControlClient = fcs;
        this.cacheKey = resultCacheKey;
        this.controlPoint = controlPoint;
        this.error = error;
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
     * Returns whether the flow should be allowed to run, based on flow ramp mode configuration and
     * Aperture Agent response. By default, flow will be allowed to run if Aperture Agent is
     * unreachable. To change this behavior, use rampMode parameter in {@link ApertureSDK#startFlow
     * }.
     *
     * @return Whether the flow should be allowed to run
     */
    public boolean shouldRun() {
        return getDecision() == FlowDecision.Accepted
                || (getDecision() == FlowDecision.Unreachable && !this.rampMode);
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
                case REJECT_REASON_NOT_SAMPLED:
                    return HttpStatus.SC_FORBIDDEN;
                case REJECT_REASON_NO_MATCHING_RAMP:
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
     * Set the result cache entry for the flow.
     *
     * @param value entry value
     * @param ttl time-to-live of the entry
     * @return upsert grpc response
     */
    public KeyUpsertResponse setResultCache(byte[] value, Duration ttl) {
        if (this.cacheKey == null) {
            return new KeyUpsertResponse(new IllegalArgumentException("Cache key not set"));
        }
        com.google.protobuf.Duration ttl_duration =
                com.google.protobuf.Duration.newBuilder().setSeconds(ttl.getSeconds()).build();
        CacheEntry entry =
                CacheEntry.newBuilder()
                        .setTtl(ttl_duration)
                        .setValue(ByteString.copyFrom(value))
                        .setKey(this.cacheKey)
                        .build();
        CacheUpsertRequest cacheUpsertRequest =
                CacheUpsertRequest.newBuilder()
                        .setControlPoint(this.controlPoint)
                        .setResultCacheEntry(entry)
                        .build();

        CacheUpsertResponse res;
        try {
            res = this.flowControlClient.cacheUpsert(cacheUpsertRequest);
        } catch (Exception e) {
            logger.debug("Aperture gRPC call failed", e);
            return new KeyUpsertResponse(e);
        }

        if (!res.hasResultCacheResponse()) {
            return new KeyUpsertResponse(new IllegalArgumentException("No cache upsert response"));
        }

        return new KeyUpsertResponse(
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheError(
                        res.getResultCacheResponse().getError()));
    }

    /**
     * Delete the result cache entry for the flow.
     *
     * @return delete grpc response
     */
    public KeyDeleteResponse deleteResultCache() {
        if (this.cacheKey == null) {
            return new KeyDeleteResponse(new IllegalArgumentException("Cache key not set"));
        }

        CacheDeleteRequest cacheDeleteRequest =
                CacheDeleteRequest.newBuilder()
                        .setControlPoint(this.controlPoint)
                        .setResultCacheKey(this.cacheKey)
                        .build();
        CacheDeleteResponse res;
        try {
            res = this.flowControlClient.cacheDelete(cacheDeleteRequest);
        } catch (Exception e) {
            logger.debug("Aperture gRPC call failed", e);
            return new KeyDeleteResponse(e);
        }

        if (!res.hasResultCacheResponse()) {
            return new KeyDeleteResponse(new IllegalArgumentException("No cache upsert response"));
        }

        return new KeyDeleteResponse(
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheError(
                        res.getResultCacheResponse().getError()));
    }

    /**
     * Retrieve the result cache entry for the flow.
     *
     * @return cache entry for the flow
     */
    public KeyLookupResponse resultCache() {
        if (this.error != null) {
            return new KeyLookupResponse(null, LookupStatus.MISS, this.error);
        }

        if (this.checkResponse == null) {
            return new KeyLookupResponse(
                    null,
                    LookupStatus.MISS,
                    new IllegalArgumentException("No cache lookup response"));
        }

        if (!this.shouldRun()) {
            return new KeyLookupResponse(
                    null, LookupStatus.MISS, new IllegalStateException("Flow was rejected"));
        }

        if (!this.checkResponse.hasCacheLookupResponse()
                || !this.checkResponse.getCacheLookupResponse().hasResultCacheResponse()) {
            return new KeyLookupResponse(
                    null,
                    LookupStatus.MISS,
                    new IllegalArgumentException("No cache lookup response"));
        }

        com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyLookupResponse lookupResponse =
                this.checkResponse.getCacheLookupResponse().getResultCacheResponse();

        return new KeyLookupResponse(
                lookupResponse.getValue(),
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheLookupStatus(
                        lookupResponse.getLookupStatus()),
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheError(
                        lookupResponse.getError()));
    }

    /**
     * Set the global cache entry for the given key.
     *
     * @param key entry key
     * @param value entry value
     * @param ttl time-to-live of the entry
     * @return upsert grpc response
     */
    public KeyUpsertResponse setGlobalCache(String key, byte[] value, Duration ttl) {
        com.google.protobuf.Duration ttl_duration =
                com.google.protobuf.Duration.newBuilder().setSeconds(ttl.getSeconds()).build();
        CacheEntry entry =
                CacheEntry.newBuilder()
                        .setTtl(ttl_duration)
                        .setValue(ByteString.copyFrom(value))
                        .build();
        CacheUpsertRequest cacheUpsertRequest =
                CacheUpsertRequest.newBuilder().putGlobalCacheEntries(key, entry).build();

        CacheUpsertResponse res;
        try {
            res = this.flowControlClient.cacheUpsert(cacheUpsertRequest);
        } catch (Exception e) {
            logger.debug("Aperture gRPC call failed", e);
            return new KeyUpsertResponse(e);
        }

        if (res.getGlobalCacheResponsesCount() == 0) {
            return new KeyUpsertResponse(new IllegalArgumentException("No cache upsert responses"));
        }

        if (!res.containsGlobalCacheResponses(key)) {
            return new KeyUpsertResponse(
                    new IllegalArgumentException("Key missing from global cache response"));
        }

        return new KeyUpsertResponse(
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheError(
                        res.getGlobalCacheResponsesOrThrow(key).getError()));
    }

    /**
     * Delete the global cache entry for the given key.
     *
     * @param key entry key
     * @return delete grpc response
     */
    public KeyDeleteResponse deleteGlobalCache(String key) {
        CacheDeleteRequest cacheDeleteRequest =
                CacheDeleteRequest.newBuilder().addGlobalCacheKeys(key).build();
        CacheDeleteResponse res;
        try {
            res = this.flowControlClient.cacheDelete(cacheDeleteRequest);
        } catch (Exception e) {
            logger.debug("Aperture gRPC call failed", e);
            return new KeyDeleteResponse(e);
        }

        if (res.getGlobalCacheResponsesCount() == 0) {
            return new KeyDeleteResponse(new IllegalArgumentException("No cache upsert response"));
        }

        if (!res.containsGlobalCacheResponses(key)) {
            return new KeyDeleteResponse(
                    new IllegalArgumentException("Key missing from global cache response"));
        }

        return new KeyDeleteResponse(
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheError(
                        res.getGlobalCacheResponsesOrThrow(key).getError()));
    }

    /**
     * Retrieve the global cache entry for the given key.
     *
     * @param key entry key
     * @return cache entry for the flow
     */
    public KeyLookupResponse globalCache(String key) {
        if (this.error != null) {
            return new KeyLookupResponse(null, LookupStatus.MISS, this.error);
        }

        if (this.checkResponse == null) {
            return new KeyLookupResponse(
                    null,
                    LookupStatus.MISS,
                    new IllegalArgumentException("No cache lookup response"));
        }

        if (!this.shouldRun()) {
            return new KeyLookupResponse(
                    null, LookupStatus.MISS, new IllegalStateException("Flow was rejected"));
        }

        if (!this.checkResponse.hasCacheLookupResponse()
                || this.checkResponse.getCacheLookupResponse().getGlobalCacheResponsesCount()
                        == 0) {
            return new KeyLookupResponse(
                    null,
                    LookupStatus.MISS,
                    new IllegalArgumentException("No cache lookup response"));
        }

        if (!this.checkResponse.getCacheLookupResponse().containsGlobalCacheResponses(key)) {
            return new KeyLookupResponse(
                    null,
                    LookupStatus.MISS,
                    new IllegalArgumentException("Key missing from global cache response"));
        }

        com.fluxninja.generated.aperture.flowcontrol.check.v1.KeyLookupResponse lookupResponse =
                this.checkResponse.getCacheLookupResponse().getGlobalCacheResponsesOrThrow(key);

        return new KeyLookupResponse(
                lookupResponse.getValue(),
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheLookupStatus(
                        lookupResponse.getLookupStatus()),
                com.fluxninja.aperture.sdk.cache.Utils.convertCacheError(
                        lookupResponse.getError()));
    }

    public Exception getError() {
        return this.error;
    }

    /**
     * Ends the flow, notifying the Aperture Agent whether it succeeded. Flow's Status is assumed to
     * be "OK" and can be set using {@link #setStatus}.
     */
    public EndResponse end() {
        if (this.ended) {
            logger.warn("Trying to end an already ended flow with status " + this.flowStatus);
            return new EndResponse(null, new IllegalStateException("Flow already ended"));
        }

        if (this.checkResponse == null) {
            logger.warn("Trying to end a flow without a check response");
            return new EndResponse(null, new IllegalStateException("Flow without check response"));
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

        ArrayList<InflightRequestRef> inflightRequestRef = new ArrayList<InflightRequestRef>();

        for (LimiterDecision decision : this.checkResponse.getLimiterDecisionsList()) {
            if (decision.getConcurrencyLimiterInfo() != null) {
                if (decision.getConcurrencyLimiterInfo().getRequestId().isEmpty()) {
                    continue;
                }

                InflightRequestRef.Builder refBuilder =
                        InflightRequestRef.newBuilder()
                                .setPolicyName(decision.getPolicyName())
                                .setPolicyHash(decision.getPolicyHash())
                                .setComponentId(decision.getComponentId())
                                .setLabel(decision.getConcurrencyLimiterInfo().getLabel())
                                .setRequestId(decision.getConcurrencyLimiterInfo().getRequestId());

                if (decision.getConcurrencyLimiterInfo().getTokensInfo() != null) {
                    refBuilder.setTokens(
                            decision.getConcurrencyLimiterInfo().getTokensInfo().getConsumed());
                }
                inflightRequestRef.add(refBuilder.build());
            }

            if (decision.getConcurrencySchedulerInfo() != null) {
                if (decision.getConcurrencySchedulerInfo().getRequestId().isEmpty()) {
                    continue;
                }

                InflightRequestRef.Builder refBuilder =
                        InflightRequestRef.newBuilder()
                                .setPolicyName(decision.getPolicyName())
                                .setPolicyHash(decision.getPolicyHash())
                                .setComponentId(decision.getComponentId())
                                .setLabel(decision.getConcurrencySchedulerInfo().getLabel())
                                .setRequestId(
                                        decision.getConcurrencySchedulerInfo().getRequestId());

                if (decision.getConcurrencySchedulerInfo().getTokensInfo() != null) {
                    refBuilder.setTokens(
                            decision.getConcurrencySchedulerInfo().getTokensInfo().getConsumed());
                }
                inflightRequestRef.add(refBuilder.build());
            }
        }

        if (inflightRequestRef.size() > 0) {
            FlowEndRequest flowEndRequest =
                    FlowEndRequest.newBuilder()
                            .setControlPoint(this.controlPoint)
                            .addAllInflightRequests(inflightRequestRef)
                            .build();

            FlowEndResponse res;
            try {
                res = this.flowControlClient.flowEnd(flowEndRequest);
            } catch (Exception e) {
                logger.debug("Aperture gRPC call failed", e);
                return new EndResponse(null, e);
            }

            return new EndResponse(res, null);
        }

        return new EndResponse(null, null);
    }
}
