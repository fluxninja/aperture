package com.fluxninja.aperture.instrumentation;

import com.fluxninja.aperture.sdk.ApertureSDK;

public class ApertureSDKWrapper {
    public ApertureSDK apertureSDK;
    public String controlPointName;
    public boolean failOpen;

    public ApertureSDKWrapper(ApertureSDK apertureSDK, String controlPointName, boolean failOpen) {
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
        this.failOpen = failOpen;
    }
}
