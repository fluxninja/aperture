package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.HttpService;

/** A builder for configuring an {@link ApertureHTTPService}. */
public class ApertureHTTPServiceBuilder {
    ApertureSDK apertureSDK;

    public ApertureHTTPServiceBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    public ApertureHTTPService build(HttpService delegate) {
        return new ApertureHTTPService(delegate, apertureSDK);
    }
}
