package com.fluxninja.aperture.instrumentation;

import com.fluxninja.aperture.sdk.ApertureSDK;

public class ApertureSDKWrapper {
    public ApertureSDK apertureSDK;
    public String controlPointName;
    public boolean rampMode;

    public ApertureSDKWrapper(ApertureSDK apertureSDK, String controlPointName, boolean rampMode) {
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
        this.rampMode = rampMode;
    }
}
