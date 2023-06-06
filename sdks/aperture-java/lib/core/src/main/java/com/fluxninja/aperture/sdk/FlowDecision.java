package com.fluxninja.aperture.sdk;

/**
 * Represents a decision made by Aperture Agent on a {@link Flow} or {@link TrafficFlow}, or a lack
 * thereof.
 */
public enum FlowDecision {
    Accepted,
    Rejected,
    Unreachable
}
