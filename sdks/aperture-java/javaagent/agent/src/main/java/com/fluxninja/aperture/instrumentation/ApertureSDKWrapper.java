package com.fluxninja.aperture.instrumentation;

import com.fluxninja.aperture.sdk.ApertureSDK;
import java.time.Duration;

public class ApertureSDKWrapper {
    public ApertureSDK apertureSDK;
    public String controlPointName;
    public boolean rampMode;
    public Duration flowTimeout;

    public ApertureSDKWrapper(
            ApertureSDK apertureSDK,
            String controlPointName,
            boolean rampMode,
            Duration flowTimeout) {
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
        this.rampMode = rampMode;
        this.flowTimeout = flowTimeout;
    }
}
