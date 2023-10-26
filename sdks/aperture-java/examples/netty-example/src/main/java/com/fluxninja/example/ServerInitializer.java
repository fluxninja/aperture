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
    boolean rampMode;
    Duration flowTimeout;
    String controlPointName;
    boolean insecureGrpc;
    String rootCertFile;

    public ServerInitializer(
            String agentAddress,
            boolean rampMode,
            Duration flowTimeout,
            String controlPointName,
            boolean insecureGrpc,
            String rootCertFile) {
        this.agentAddress = agentAddress;
        this.rampMode = rampMode;
        this.flowTimeout = flowTimeout;
        this.controlPointName = controlPointName;
        this.insecureGrpc = insecureGrpc;
        this.rootCertFile = rootCertFile;
    }

    @Override
    protected void initChannel(Channel ch) {
        try {
            sdk =
                    ApertureSDK.builder()
                            .setAddress(this.agentAddress)
                            .useInsecureGrpc(insecureGrpc)
                            .setRootCertificateFile(rootCertFile)
                            .build();
        } catch (IOException ex) {
            throw new RuntimeException(ex);
        }

        ChannelPipeline pipeline = ch.pipeline();
        pipeline.addLast(new HttpServerCodec());
        pipeline.addLast(new HttpObjectAggregator(Integer.MAX_VALUE));
        // ApertureServerHandler must be added before the response-generating HelloWorldHandler,
        //    but after the codec handler.
        pipeline.addLast(new ApertureServerHandler(sdk, controlPointName, rampMode, flowTimeout));
        pipeline.addLast(new HelloWorldHandler());
    }
}
