package com.fluxninja.aperture.sdk.cache;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheLookupStatus;

public class Utils {
    public static LookupStatus convertCacheLookupStatus(CacheLookupStatus status) {
        return (status == CacheLookupStatus.HIT) ? LookupStatus.HIT : LookupStatus.MISS;
    }

    public static Exception convertCacheError(String errorMessage) {
        return (errorMessage == null || errorMessage.isEmpty())
                ? null
                : new Exception(errorMessage);
    }
}
