package com.fluxninja.aperture.netty;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.Constants;
import com.fluxninja.aperture.sdk.EndResponse;
import com.fluxninja.aperture.sdk.FlowDecision;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.aperture.sdk.TrafficFlowRequest;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpHeaderNames;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import io.netty.util.CharsetUtil;
import java.time.Duration;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

public class ApertureServerHandler extends SimpleChannelInboundHandler<HttpRequest> {

    private final ApertureSDK apertureSDK;
    private final String controlPointName;
    private boolean rampMode = false;
    private Duration flowTimeout = Constants.DEFAULT_RPC_TIMEOUT;

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

    public ApertureServerHandler(
            ApertureSDK sdk, String controlPointName, boolean rampMode, Duration flowTimeout) {
        if (controlPointName == null || controlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must not be null or empty");
        }
        if (sdk == null) {
            throw new IllegalArgumentException("Aperture SDK must not be null");
        }
        this.apertureSDK = sdk;
        this.controlPointName = controlPointName;
        this.rampMode = rampMode;
        this.flowTimeout = flowTimeout;
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, HttpRequest req) {
        TrafficFlowRequest trafficFlowRequest =
                NettyUtils.trafficFlowRequestFromRequest(ctx, req, controlPointName, flowTimeout);

        TrafficFlow flow = this.apertureSDK.startTrafficFlow(trafficFlowRequest);

        if (flow.ignored()) {
            ctx.fireChannelRead(req);
            return;
        }

        FlowDecision flowDecision = flow.getDecision();
        boolean flowAccepted =
                (flowDecision == FlowDecision.Accepted
                        || (flowDecision == FlowDecision.Unreachable && !this.rampMode));

        if (flowAccepted) {
            try {
                Map<String, String> newHeaders = new HashMap<>();
                if (flow.checkResponse() != null) {
                    newHeaders = flow.checkResponse().getOkResponse().getHeadersMap();
                }
                HttpRequest newRequest = NettyUtils.updateHeaders(req, newHeaders);

                ctx.fireChannelRead(newRequest);
            } catch (Exception e) {
                flow.setStatus(FlowStatus.Error);
                throw e;
            } finally {
                EndResponse endResponse = flow.end();
                if (endResponse.getError() != null) {
                    System.err.println("Error ending flow: " + endResponse.getError().getMessage());
                }
            }
        } else {
            flow.setStatus(FlowStatus.Unset);
            EndResponse endResponse = flow.end();
            if (endResponse.getError() != null) {
                System.err.println("Error ending flow: " + endResponse.getError().getMessage());
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
