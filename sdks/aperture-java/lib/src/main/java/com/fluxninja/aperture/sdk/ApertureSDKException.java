package com.fluxninja.aperture.sdk;

public class ApertureSDKException extends Exception {
    public ApertureSDKException(String message) {
        super(message);
    }

    public ApertureSDKException(Exception e) {
        super(e.getMessage(), e);
    }
}
