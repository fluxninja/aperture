package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.HttpService;

/** A builder for configuring an {@link ApertureHTTPService}. */
public class ApertureHTTPServiceBuilder {
    private ApertureSDK apertureSDK;
    private String controlPointName;

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

    public ApertureHTTPService build(HttpService delegate) {
        if (this.controlPointName == null || this.controlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must be set");
        }
        if (this.apertureSDK == null) {
            throw new IllegalArgumentException("Aperture SDK must be set");
        }
        return new ApertureHTTPService(delegate, apertureSDK, controlPointName);
    }
}
