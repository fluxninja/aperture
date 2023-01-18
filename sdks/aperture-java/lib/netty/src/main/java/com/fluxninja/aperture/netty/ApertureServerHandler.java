package com.fluxninja.aperture.netty;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundInvoker;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.*;
import io.netty.util.CharsetUtil;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class ApertureServerHandler extends SimpleChannelInboundHandler<HttpRequest> {

    private final ApertureSDK apertureSDK;

    public ApertureServerHandler(ApertureSDK sdk) {
        this.apertureSDK = sdk;
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, HttpRequest req) {

        AttributeContext attributes = NettyUtils.attributesFromRequest(req);
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(attributes);

        if (flow.accepted()) {
            try {
                List<HeaderValueOption> newHeaders = flow.checkResponse().getOkResponse().getHeadersList();
                HttpRequest newRequest = NettyUtils.updateHeaders(req, newHeaders);

                ctx.fireChannelRead(newRequest);
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                // ending flow failed
                e.printStackTrace();
                FullHttpResponse response = new DefaultFullHttpResponse(HttpVersion.HTTP_1_1, HttpResponseStatus.INTERNAL_SERVER_ERROR);
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
            if (flow.checkResponse().hasDeniedResponse() && flow.checkResponse().getDeniedResponse().hasStatus()) {
                status = HttpResponseStatus.valueOf(flow.checkResponse().getDeniedResponse().getStatus().getCodeValue());
            } else {
                status = HttpResponseStatus.FORBIDDEN;
            }

            ByteBuf content = Unpooled.copiedBuffer(status.toString(), CharsetUtil.UTF_8);
            FullHttpResponse response = new DefaultFullHttpResponse(HttpVersion.HTTP_1_1, status, content);
            response.headers().set(HttpHeaderNames.CONTENT_TYPE, "text/html");
            response.headers().set(HttpHeaderNames.CONTENT_LENGTH, content.readableBytes());
            ctx.write(response);
            ctx.flush();
        }
    }
}
