package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.envoy.service.auth.v3.Address;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.SocketAddress;
import java.time.Instant;

public class Utils {
    public static long getCurrentEpochNanos() {
        long nanosInSecond = 1000000000L;
        Instant currentTime = Instant.now();
        return currentTime.getEpochSecond() * nanosInSecond + currentTime.getNano();
    }

    public static AttributeContext.Peer peerFromAddress(String address) {
        return AttributeContext.Peer.newBuilder()
                .setAddress(
                        Address.newBuilder()
                                .setSocketAddress(
                                        SocketAddress.newBuilder().setAddress(address).build())
                                .build())
                .build();
    }
}
