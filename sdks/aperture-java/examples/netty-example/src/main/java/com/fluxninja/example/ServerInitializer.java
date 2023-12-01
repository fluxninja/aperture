package com.fluxninja.example;

import com.fluxninja.aperture.netty.ApertureServerHandler;
import com.fluxninja.aperture.sdk.ApertureSDK;
import io.netty.channel.Channel;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.codec.http.HttpServerCodec;
import java.io.IOException;
import java.time.Duration;

public class ServerInitializer extends ChannelInitializer<Channel> {

    ApertureSDK sdk;
    String agentAddress;
    String agentAPIKey;
    boolean rampMode;
    Duration flowTimeout;
    String controlPointName;
    boolean insecureGrpc;
    String rootCertFile;

    public ServerInitializer(
            String agentAddress,
            String agentAPIKey,
            boolean rampMode,
            Duration flowTimeout,
            String controlPointName,
            boolean insecureGrpc,
            String rootCertFile) {
        this.agentAddress = agentAddress;
        this.agentAPIKey = agentAPIKey;
        this.rampMode = rampMode;
        this.flowTimeout = flowTimeout;
        this.controlPointName = controlPointName;
        this.insecureGrpc = insecureGrpc;
        this.rootCertFile = rootCertFile;
    }

    // START: NettyInitChannel

    @Override
    protected void initChannel(Channel ch) {

        // START: NettyCreateSDK
        try {
            sdk =
                    ApertureSDK.builder()
                            .setAddress(this.agentAddress)
                            .setAPIKey(this.agentAPIKey)
                            .useInsecureGrpc(insecureGrpc) // Optional: Defaults to true
                            .setRootCertificateFile(rootCertFile)
                            .addIgnoredPaths("/health,/connected")
                            .build();
        } catch (IOException ex) {
            throw new RuntimeException(ex);
        }
        // END: NettyCreateSDK

        ChannelPipeline pipeline = ch.pipeline();
        pipeline.addLast(new HttpServerCodec());
        pipeline.addLast(new HttpObjectAggregator(Integer.MAX_VALUE));
        // ApertureServerHandler must be added before the response-generating
        // HelloWorldHandler,
        // but after the codec handler.
        pipeline.addLast(new ApertureServerHandler(sdk, controlPointName, rampMode, flowTimeout));
        pipeline.addLast(new HelloWorldHandler());
    }
    // END: NettyInitChannel
}
