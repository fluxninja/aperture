package com.fluxninja.example;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import java.time.Duration;

public class NettyServer {

    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_ADDRESS = "localhost:8089";
    public static final String DEFAULT_RAMP_MODE = "false";
    public static final String DEFAULT_CONTROL_POINT_NAME = "awesome_feature";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

    public static void main(String[] args) throws Exception {
        String agentAddress = System.getenv("APERTURE_AGENT_ADDRESS");
        if (agentAddress == null) {
            agentAddress = DEFAULT_AGENT_ADDRESS;
        }
        String agentAPIKey = System.getenv("APERTURE_AGENT_API_KEY");
        if (agentAPIKey == null) {
            agentAPIKey = "";
        }
        String appPort = System.getenv("APERTURE_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }
        String rampModeString = System.getenv("APERTURE_ENABLE_RAMP_MODE");
        if (rampModeString == null) {
            rampModeString = DEFAULT_RAMP_MODE;
        }
        boolean rampMode = Boolean.parseBoolean(rampModeString);

        String controlPointName = System.getenv("APERTURE_CONTROL_POINT_NAME");
        if (controlPointName == null) {
            controlPointName = DEFAULT_CONTROL_POINT_NAME;
        }
        String insecureGrpcString = System.getenv("APERTURE_AGENT_INSECURE");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String rootCertFile = System.getenv("APERTURE_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }

        Duration flowTimeout = Duration.ofMillis(1000);

        EventLoopGroup bossGroup = new NioEventLoopGroup();
        EventLoopGroup workerGroup = new NioEventLoopGroup();

        try {

            ServerBootstrap httpBootstrap = new ServerBootstrap();

            httpBootstrap
                    .group(bossGroup, workerGroup)
                    .channel(NioServerSocketChannel.class)
                    .childHandler(
                            new ServerInitializer(
                                    agentAddress,
                                    agentAPIKey,
                                    rampMode,
                                    flowTimeout,
                                    controlPointName,
                                    insecureGrpc,
                                    rootCertFile))
                    .option(ChannelOption.SO_BACKLOG, 128)
                    .childOption(ChannelOption.SO_KEEPALIVE, true);

            // Bind and start to accept incoming connections.
            ChannelFuture httpChannel = httpBootstrap.bind(Integer.parseInt(appPort)).sync();

            // Wait until the server socket is closed
            httpChannel.channel().closeFuture().sync();
        } finally {
            workerGroup.shutdownGracefully();
            bossGroup.shutdownGracefully();
        }
    }
}
