package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.client.HttpClient;

/** A builder for configuring an {@link ApertureHTTPClient}. */
public class ApertureHTTPClientBuilder {
    private ApertureSDK apertureSDK;
    private String controlPointName = "";

    public ApertureHTTPClientBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    /**
     * Sets the control point name for traffic produced by this client.
     *
     * @param controlPointName control point name to be used
     * @return the builder object.
     */
    public ApertureHTTPClientBuilder setControlPointName(String controlPointName) {
        this.controlPointName = controlPointName;
        return this;
    }

    public ApertureHTTPClient build(HttpClient delegate) {
        return new ApertureHTTPClient(delegate, apertureSDK, controlPointName);
    }
}
