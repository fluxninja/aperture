package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.client.HttpClient;

/** A builder for configuring an {@link ApertureHTTPClient}. */
public class ApertureHTTPClientBuilder {
    ApertureSDK apertureSDK;
    String controlPointName = "egress";

    public ApertureHTTPClientBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    public ApertureHTTPClientBuilder setControlPointName(String controlPointName) {
        this.controlPointName = controlPointName;
        return this;
    }

    public ApertureHTTPClient build(HttpClient delegate) {
        return new ApertureHTTPClient(delegate, apertureSDK, controlPointName);
    }
}
