package com.fluxninja.aperture.sdk.cache;

public class KeyDeleteResponse {
    private final Exception error;

    public KeyDeleteResponse(Exception error) {
        this.error = error;
    }

    public Exception getError() {
        return error;
    }
}
