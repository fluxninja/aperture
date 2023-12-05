package com.fluxninja.aperture.sdk.cache;

public class KeyLookupResponse {
    private final Object value;
    private final LookupStatus lookupStatus;
    private final Exception error;

    public KeyLookupResponse(Object value, LookupStatus lookupStatus, Exception error) {
        this.value = value;
        this.lookupStatus = lookupStatus;
        this.error = error;
    }

    public Object getValue() {
        return value;
    }

    public LookupStatus getLookupStatus() {
        return lookupStatus;
    }

    public Exception getError() {
        return error;
    }
}
