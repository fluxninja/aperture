package com.fluxninja.example;

import com.fluxninja.aperture.netty.ApertureServerHandler;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import io.netty.channel.Channel;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.codec.http.HttpServerCodec;

public class ServerInitializer extends ChannelInitializer<Channel> {

    ApertureSDK sdk;
    String agentHost;
    int agentPort;


    public ServerInitializer(String agentHost, String agentPort) {
        this.agentHost = agentHost;
        this.agentPort = Integer.parseInt(agentPort);
    }

    @Override
    protected void initChannel(Channel ch) {
        try {
            sdk = ApertureSDK.builder().setHost(this.agentHost).setPort(this.agentPort).build();
        } catch (ApertureSDKException ex) {
            throw new RuntimeException(ex);
        }

        ChannelPipeline pipeline = ch.pipeline();
        pipeline.addLast(new HttpServerCodec());
        pipeline.addLast(new HttpObjectAggregator(Integer.MAX_VALUE));
        // ApertureServerHandler must be added before the response-generating HelloWorldHandler,
        //    but after the codec handler.
        pipeline.addLast(new ApertureServerHandler(sdk));
        pipeline.addLast(new HelloWorldHandler());
    }

}
