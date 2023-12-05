package com.fluxninja.aperture.sdk.cache;

public enum LookupStatus {
    HIT("HIT"),
    MISS("MISS");

    private final String value;

    LookupStatus(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }
}
