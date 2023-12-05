package com.fluxninja.aperture.sdk.cache;

public class KeyUpsertResponse {
    private final Exception error;

    public KeyUpsertResponse(Exception error) {
        this.error = error;
    }

    public Exception getError() {
        return error;
    }
}
