package com.fluxninja.aperture.netty;

import com.fluxninja.aperture.sdk.*;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.*;
import io.netty.util.CharsetUtil;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

public class ApertureServerHandler extends SimpleChannelInboundHandler<HttpRequest> {

    private final ApertureSDK apertureSDK;
    private final String controlPointName;
    private boolean failOpen = true;

    public ApertureServerHandler(ApertureSDK sdk, String controlPointName) {
        if (controlPointName == null || controlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must not be null or empty");
        }
        if (sdk == null) {
            throw new IllegalArgumentException("Aperture SDK must not be null");
        }
        this.apertureSDK = sdk;
        this.controlPointName = controlPointName;
    }

    public ApertureServerHandler(ApertureSDK sdk, String controlPointName, boolean failOpen) {
        if (controlPointName == null || controlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must not be null or empty");
        }
        if (sdk == null) {
            throw new IllegalArgumentException("Aperture SDK must not be null");
        }
        this.apertureSDK = sdk;
        this.controlPointName = controlPointName;
        this.failOpen = failOpen;
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, HttpRequest req) {
        CheckHTTPRequest checkRequest =
                NettyUtils.checkRequestFromRequest(ctx, req, controlPointName);
        String path = new QueryStringDecoder(req.uri()).path();

        TrafficFlow flow = this.apertureSDK.startTrafficFlow(path, checkRequest);

        if (flow.ignored()) {
            ctx.fireChannelRead(req);
            return;
        }

        FlowDecision flowDecision = flow.getDecision();
        boolean flowAccepted =
                (flowDecision == FlowDecision.Accepted
                        || (flowDecision == FlowDecision.Unreachable && this.failOpen));

        if (flowAccepted) {
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
            Map<String, String> headers;
            if (flow.checkResponse() != null && flow.checkResponse().hasDeniedResponse()) {
                status = HttpResponseStatus.valueOf(flow.getRejectionHttpStatusCode());
                headers = flow.checkResponse().getDeniedResponse().getHeadersMap();

            } else {
                status = HttpResponseStatus.FORBIDDEN;
                headers = Collections.emptyMap();
            }

            ByteBuf content = Unpooled.copiedBuffer(status.toString(), CharsetUtil.UTF_8);
            FullHttpResponse response =
                    new DefaultFullHttpResponse(HttpVersion.HTTP_1_1, status, content);

            for (Map.Entry<String, String> entry : headers.entrySet()) {
                response.headers().set(entry.getKey(), entry.getValue());
            }
            response.headers().set(HttpHeaderNames.CONTENT_TYPE, "text/html");
            response.headers().set(HttpHeaderNames.CONTENT_LENGTH, content.readableBytes());

            ctx.write(response);
            ctx.flush();
        }
    }
}
