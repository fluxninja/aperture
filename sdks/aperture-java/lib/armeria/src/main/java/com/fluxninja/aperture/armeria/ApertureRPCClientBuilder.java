package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.client.RpcClient;

/** A builder for configuring an {@link ApertureRPCClient}. */
public class ApertureRPCClientBuilder {
    ApertureSDK apertureSDK;

    public ApertureRPCClientBuilder setApertureSDK(ApertureSDK apertureSDK) {
        this.apertureSDK = apertureSDK;
        return this;
    }

    public ApertureRPCClient build(RpcClient delegate) {
        return new ApertureRPCClient(delegate, apertureSDK);
    }
}
