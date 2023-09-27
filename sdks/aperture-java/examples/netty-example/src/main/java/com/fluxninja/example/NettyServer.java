package com.fluxninja.example;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;

public class NettyServer {

    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";
    public static final String DEFAULT_RAMP_MODE = "false";
    public static final String DEFAULT_CONTROL_POINT_NAME = "awesome_feature";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

    public static void main(String[] args) throws Exception {
        String agentHost = System.getenv("APERTURE_AGENT_HOST");
        if (agentHost == null) {
            agentHost = DEFAULT_AGENT_HOST;
        }
        String agentPort = System.getenv("APERTURE_AGENT_PORT");
        if (agentPort == null) {
            agentPort = DEFAULT_AGENT_PORT;
        }
        String appPort = System.getenv("FN_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }
        String rampModeString = System.getenv("FN_ENABLE_RAMP_MODE");
        if (rampModeString == null) {
            rampModeString = DEFAULT_RAMP_MODE;
        }
        boolean rampMode = Boolean.parseBoolean(rampModeString);

        String controlPointName = System.getenv("FN_CONTROL_POINT_NAME");
        if (controlPointName == null) {
            controlPointName = DEFAULT_CONTROL_POINT_NAME;
        }
        String insecureGrpcString = System.getenv("FN_INSECURE_GRPC");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String rootCertFile = System.getenv("FN_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }

        EventLoopGroup bossGroup = new NioEventLoopGroup();
        EventLoopGroup workerGroup = new NioEventLoopGroup();

        try {

            ServerBootstrap httpBootstrap = new ServerBootstrap();

            httpBootstrap
                    .group(bossGroup, workerGroup)
                    .channel(NioServerSocketChannel.class)
                    .childHandler(
                            new ServerInitializer(
                                    agentHost,
                                    agentPort,
                                    rampMode,
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
