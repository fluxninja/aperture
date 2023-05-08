package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.client.RpcClient;

/** A builder for configuring an {@link ApertureRPCClient}. */
public class ApertureRPCClientBuilder {
    private ApertureSDK apertureSDK;
    private String controlPointName;
    private boolean enableFailOpen = true;

    /**
     * Sets the Aperture SDK used by this service.
     *
     * @param apertureSDK instance of Aperture SDK to be used
     * @return the builder object.
     */
    public ApertureRPCClientBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    /**
     * Sets the control point name for traffic produced by this client.
     *
     * @param controlPointName control point name to be used
     * @return the builder object.
     */
    public ApertureRPCClientBuilder setControlPointName(String controlPointName) {
        this.controlPointName = controlPointName;
        return this;
    }

    /**
     * Sets the fail-open behavior for the client when the Aperture Agent is unreachable. If set to
     * true, all traffic will pass through; if set to false, all traffic will be blocked.
     *
     * @param enableFailOpen whether all traffic should be accepted when Aperture Agent is
     *     unreachable
     * @return the builder object.
     */
    public ApertureRPCClientBuilder setEnableFailOpen(boolean enableFailOpen) {
        this.enableFailOpen = enableFailOpen;
        return this;
    }

    public ApertureRPCClient build(RpcClient delegate) {
        if (this.controlPointName == null || this.controlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must be set");
        }
        if (this.apertureSDK == null) {
            throw new IllegalArgumentException("Aperture SDK must be set");
        }
        return new ApertureRPCClient(delegate, apertureSDK, controlPointName, enableFailOpen);
    }
}
