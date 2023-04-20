package com.fluxninja.example;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;

public class NettyServer {

    public static final String DEFAULT_APP_PORT = "8080";

    public static void main(String[] args) throws Exception {
        String appPort = System.getenv("FN_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }

        EventLoopGroup bossGroup = new NioEventLoopGroup();
        EventLoopGroup workerGroup = new NioEventLoopGroup();

        try {

            ServerBootstrap httpBootstrap = new ServerBootstrap();

            httpBootstrap
                    .group(bossGroup, workerGroup)
                    .channel(NioServerSocketChannel.class)
                    .childHandler(new ServerInitializer())
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
