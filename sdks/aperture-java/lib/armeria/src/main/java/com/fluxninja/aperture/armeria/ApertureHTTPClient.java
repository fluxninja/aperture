package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.*;
import com.linecorp.armeria.client.ClientRequestContext;
import com.linecorp.armeria.client.HttpClient;
import com.linecorp.armeria.client.SimpleDecoratingHttpClient;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpResponse;
import java.util.Collections;
import java.util.Map;
import java.util.function.Function;

/** Decorates an {@link HttpClient} to enable flow control using provided {@link ApertureSDK} */
public class ApertureHTTPClient extends SimpleDecoratingHttpClient {
    private final ApertureSDK apertureSDK;
    private final String controlPointName;
    private final boolean rampMode;

    public static Function<? super HttpClient, ApertureHTTPClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName) {
        ApertureHTTPClientBuilder builder = new ApertureHTTPClientBuilder();
        builder.setApertureSDK(apertureSDK).setControlPointName(controlPointName);
        return builder::build;
    }

    public static Function<? super HttpClient, ApertureHTTPClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName, boolean rampMode) {
        ApertureHTTPClientBuilder builder = new ApertureHTTPClientBuilder();
        builder.setApertureSDK(apertureSDK)
                .setControlPointName(controlPointName)
                .setEnableRampMode(rampMode);
        return builder::build;
    }

    public ApertureHTTPClient(
            HttpClient delegate,
            ApertureSDK apertureSDK,
            String controlPointName,
            boolean rampMode) {
        super(delegate);
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
        this.rampMode = rampMode;
    }

    @Override
    public HttpResponse execute(ClientRequestContext ctx, HttpRequest req) throws Exception {
        TrafficFlowRequest request =
                HttpUtils.trafficFlowRequestFromRequest(ctx, req, this.controlPointName);
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(request);

        if (flow.ignored()) {
            return unwrap().execute(ctx, req);
        }

        FlowDecision flowDecision = flow.getDecision();
        boolean flowAccepted =
                (flowDecision == FlowDecision.Accepted
                        || (flowDecision == FlowDecision.Unreachable && !this.rampMode));

        if (flowAccepted) {
            HttpResponse res;
            try {
                Map<String, String> newHeaders = Collections.emptyMap();
                if (flow.checkResponse() != null) {
                    newHeaders = flow.checkResponse().getOkResponse().getHeadersMap();
                }
                HttpRequest newRequest = HttpUtils.updateHeaders(req, newHeaders);
                ctx.updateRequest(newRequest);

                res = unwrap().execute(ctx, newRequest);
            } catch (Exception e) {
                flow.setStatus(FlowStatus.Error);
                throw e;
            } finally {
                flow.end();
            }
            return res;
        } else {
            return HttpUtils.handleRejectedFlow(flow);
        }
    }
}
