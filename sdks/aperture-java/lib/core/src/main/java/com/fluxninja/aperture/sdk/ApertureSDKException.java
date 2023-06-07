package com.fluxninja.aperture.sdk;

/** Exception thrown by Aperture SDK when it cannot be constructed, or when used incorrectly. */
public class ApertureSDKException extends Exception {
    public ApertureSDKException(String message) {
        super(message);
    }

    public ApertureSDKException(Exception e) {
        super(e.getMessage(), e);
    }
}
