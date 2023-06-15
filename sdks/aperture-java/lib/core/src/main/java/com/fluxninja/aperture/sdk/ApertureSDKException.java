package com.fluxninja.aperture.sdk;

/** Exception thrown by Aperture SDK when it cannot be constructed. */
public class ApertureSDKException extends Exception {
    public ApertureSDKException(String message) {
        super(message);
    }

    public ApertureSDKException(Exception e) {
        super(e.getMessage(), e);
    }
}
