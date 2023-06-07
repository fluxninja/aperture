package com.fluxninja.aperture.sdk;

/**
 * Status of a finished flow, sent to Aperture Agent when calling {@link
 * com.fluxninja.aperture.sdk.Flow#end Flow#end}
 */
public enum FlowStatus {
    OK,
    Error,
    Unset,
}
