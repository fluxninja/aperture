package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import java.time.Duration;

/**
 * The TrafficFlowRequest class represents the parameters for initiating a traffic flow within the
 * Aperture SDK. For more details, see ApertureSDK's {@link ApertureSDK#startTrafficFlow
 * startTrafficFlow} method.
 */
public class TrafficFlowRequest {
    private final CheckHTTPRequest checkHTTPRequest;
    private final Duration flowTimeout;

    /**
     * Constructs a new TrafficFlowRequest object with the specified CheckHTTPRequest.
     *
     * @param checkHTTPRequest The CheckHTTPRequest object to be encapsulated.
     */
    TrafficFlowRequest(CheckHTTPRequest checkHTTPRequest, Duration flowTimeout) {
        this.checkHTTPRequest = checkHTTPRequest;
        this.flowTimeout = flowTimeout;
    }

    public Duration getFlowTimeout() {
        return flowTimeout;
    }

    /**
     * Returns the encapsulated CheckHTTPRequest object.
     *
     * @return The CheckHTTPRequest object.
     */
    public CheckHTTPRequest getCheckHTTPRequest() {
        return checkHTTPRequest;
    }

    /**
     * Creates a new instance of TrafficFlowRequestBuilder for building TrafficFlowRequest objects.
     *
     * @return A new TrafficFlowRequestBuilder object.
     */
    public static TrafficFlowRequestBuilder newBuilder() {
        return new TrafficFlowRequestBuilder();
    }
}
