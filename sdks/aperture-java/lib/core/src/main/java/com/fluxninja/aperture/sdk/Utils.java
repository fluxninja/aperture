package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.SocketAddress;
import java.time.Instant;

public class Utils {
    public static long getCurrentEpochNanos() {
        long nanosInSecond = 1000000000L;
        Instant currentTime = Instant.now();
        return currentTime.getEpochSecond() * nanosInSecond + currentTime.getNano();
    }

    public static SocketAddress createSocketAddress(String address, int port, String protocol) {
        return SocketAddress.newBuilder()
                .setAddress(address)
                .setPort(port)
                .setProtocol(SocketAddress.Protocol.valueOf(protocol))
                .build();
    }
}
