package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.RpcService;

/** A builder for configuring an {@link ApertureRPCService}. */
public class ApertureRPCServiceBuilder {
    ApertureSDK apertureSDK;
    public ApertureRPCServiceBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    public ApertureRPCService build(RpcService delegate) {
        return new ApertureRPCService(delegate, apertureSDK);
    }
}
