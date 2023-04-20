package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.HttpService;

/** A builder for configuring an {@link ApertureHTTPService}. */
public class ApertureHTTPServiceBuilder {
    private ApertureSDK apertureSDK;
    private String controlPointName = "ingress";

    public ApertureHTTPServiceBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    /**
     * Sets the control point name for traffic handled by this service. If not set, defaults to
     * "ingress".
     *
     * @param controlPointName control point name to be used
     * @return the builder object.
     */
    public ApertureHTTPServiceBuilder setControlPointName(String controlPointName) {
        this.controlPointName = controlPointName;
        return this;
    }

    public ApertureHTTPService build(HttpService delegate) {
        return new ApertureHTTPService(delegate, apertureSDK, controlPointName);
    }
}
