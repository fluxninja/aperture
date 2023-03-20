package com.fluxninja.aperture.sdk;

import java.time.Instant;

public class Utils {
    public static long getCurrentEpochNanos() {
        long nanosInSecond = 1000000000L;
        Instant currentTime = Instant.now();
        return currentTime.getEpochSecond() * nanosInSecond + currentTime.getNano();
    }
}
