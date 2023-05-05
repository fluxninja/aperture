package com.fluxninja.aperture.instrumentation;

import com.fluxninja.aperture.sdk.ApertureSDK;

public class ApertureSDKWrapper {
    public ApertureSDK apertureSDK;
    public String controlPointName;

    public ApertureSDKWrapper(ApertureSDK apertureSDK, String controlPointName) {
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
    }
}
