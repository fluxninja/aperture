package com.fluxninja.aperture.netty;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.*;
import io.netty.util.CharsetUtil;
import java.util.HashMap;
import java.util.Map;

public class ApertureServerHandler extends SimpleChannelInboundHandler<HttpRequest> {

    private final ApertureSDK apertureSDK;

    public ApertureServerHandler(ApertureSDK sdk) {
        this.apertureSDK = sdk;
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, HttpRequest req) {
        CheckHTTPRequest checkRequest = NettyUtils.checkRequestFromRequest(ctx, req);
        String path = new QueryStringDecoder(req.uri()).path();

        TrafficFlow flow = this.apertureSDK.startTrafficFlow(path, checkRequest);

        if (flow.ignored()) {
            ctx.fireChannelRead(req);
            return;
        }

        if (flow.accepted()) {
            try {
                Map<String, String> newHeaders = new HashMap<>();
                if (flow.checkResponse() != null) {
                    newHeaders = flow.checkResponse().getOkResponse().getHeadersMap();
                }
                HttpRequest newRequest = NettyUtils.updateHeaders(req, newHeaders);

                ctx.fireChannelRead(newRequest);
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                // ending flow failed
                e.printStackTrace();
                FullHttpResponse response =
                        new DefaultFullHttpResponse(
                                HttpVersion.HTTP_1_1, HttpResponseStatus.INTERNAL_SERVER_ERROR);
                ctx.write(response);
                ctx.flush();
            } catch (Exception e) {
                try {
                    flow.end(FlowStatus.Error);
                } catch (ApertureSDKException ae) {
                    e.printStackTrace();
                    ae.printStackTrace();
                }
                throw e;
            }
        } else {
            try {
                flow.end(FlowStatus.Unset);
            } catch (ApertureSDKException e) {
                e.printStackTrace();
            }
            HttpResponseStatus status;
            if (flow.checkResponse() != null
                    && flow.checkResponse().hasDeniedResponse()
                    && flow.checkResponse().getDeniedResponse().getStatus() != 0) {
                status =
                        HttpResponseStatus.valueOf(
                                flow.checkResponse().getDeniedResponse().getStatus());
            } else {
                status = HttpResponseStatus.FORBIDDEN;
            }

            ByteBuf content = Unpooled.copiedBuffer(status.toString(), CharsetUtil.UTF_8);
            FullHttpResponse response =
                    new DefaultFullHttpResponse(HttpVersion.HTTP_1_1, status, content);
            response.headers().set(HttpHeaderNames.CONTENT_TYPE, "text/html");
            response.headers().set(HttpHeaderNames.CONTENT_LENGTH, content.readableBytes());
            ctx.write(response);
            ctx.flush();
        }
    }
}
