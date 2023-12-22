package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.Constants;
import com.linecorp.armeria.server.HttpService;
import java.time.Duration;

/** A builder for configuring an {@link ApertureHTTPService}. */
public class ApertureHTTPServiceBuilder {
    private ApertureSDK apertureSDK;
    private String controlPointName;
    private boolean enableRampMode = false;
    private Duration flowTimeout = Constants.DEFAULT_RPC_TIMEOUT;

    /**
     * Sets timeout for connection to Aperture Agent. Set to 0 to block until response is received.
     *
     * @param flowTimeout The timeout for connection to Aperture Agent.
     * @return The builder object.
     */
    public ApertureHTTPServiceBuilder setFlowTimeout(Duration flowTimeout) {
        this.flowTimeout = flowTimeout;
        return this;
    }

    /**
     * Sets the Aperture SDK used by this service.
     *
     * @param apertureSDK instance of Aperture SDK to be used
     * @return the builder object.
     */
    public ApertureHTTPServiceBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    /**
     * Sets the control point name for traffic handled by this service.
     *
     * @param controlPointName control point name to be used
     * @return the builder object.
     */
    public ApertureHTTPServiceBuilder setControlPointName(String controlPointName) {
        this.controlPointName = controlPointName;
        return this;
    }

    /**
     * Marks started flows as ramp mode, requiring at least one ramp component to accept it. Marked
     * flows will fail if the policy is not loaded or Agent is unreachable.
     *
     * @param enableRampMode whether all started flows should be started in ramp mode
     * @return the builder object.
     */
    public ApertureHTTPServiceBuilder setEnableRampMode(boolean enableRampMode) {
        this.enableRampMode = enableRampMode;
        return this;
    }

    public ApertureHTTPService build(HttpService delegate) {
        if (this.controlPointName == null || this.controlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must be set");
        }
        if (this.apertureSDK == null) {
            throw new IllegalArgumentException("Aperture SDK must be set");
        }
        return new ApertureHTTPService(
                delegate, apertureSDK, controlPointName, enableRampMode, flowTimeout);
    }
}
