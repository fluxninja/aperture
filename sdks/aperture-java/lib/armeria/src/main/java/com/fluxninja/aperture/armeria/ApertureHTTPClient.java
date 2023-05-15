package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.*;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import com.linecorp.armeria.client.ClientRequestContext;
import com.linecorp.armeria.client.HttpClient;
import com.linecorp.armeria.client.SimpleDecoratingHttpClient;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpResponse;
import com.linecorp.armeria.common.HttpStatus;
import java.util.Collections;
import java.util.Map;
import java.util.function.Function;

/** Decorates an {@link HttpClient} to enable flow control using provided {@link ApertureSDK} */
public class ApertureHTTPClient extends SimpleDecoratingHttpClient {
    private final ApertureSDK apertureSDK;
    private final String controlPointName;
    private final boolean failOpen;

    public static Function<? super HttpClient, ApertureHTTPClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName) {
        ApertureHTTPClientBuilder builder = new ApertureHTTPClientBuilder();
        builder.setApertureSDK(apertureSDK).setControlPointName(controlPointName);
        return builder::build;
    }

    public static Function<? super HttpClient, ApertureHTTPClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName, boolean failOpen) {
        ApertureHTTPClientBuilder builder = new ApertureHTTPClientBuilder();
        builder.setApertureSDK(apertureSDK)
                .setControlPointName(controlPointName)
                .setEnableFailOpen(failOpen);
        return builder::build;
    }

    public ApertureHTTPClient(
            HttpClient delegate,
            ApertureSDK apertureSDK,
            String controlPointName,
            boolean failOpen) {
        super(delegate);
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
        this.failOpen = failOpen;
    }

    @Override
    public HttpResponse execute(ClientRequestContext ctx, HttpRequest req) throws Exception {
        CheckHTTPRequest request =
                HttpUtils.checkRequestFromRequest(ctx, req, this.controlPointName);
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(req.path(), request);

        if (flow.ignored()) {
            return unwrap().execute(ctx, req);
        }

        FlowResult flowResult = flow.result();
        boolean flowAccepted =
                (flowResult == FlowResult.Accepted
                        || (flowResult == FlowResult.Unreachable && this.failOpen));

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
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                // ending flow failed
                e.printStackTrace();
                return HttpResponse.of(HttpStatus.INTERNAL_SERVER_ERROR);
            } catch (Exception e) {
                try {
                    flow.end(FlowStatus.Error);
                } catch (ApertureSDKException ae) {
                    ae.printStackTrace();
                }
                throw e;
            }
            return res;
        } else {
            return HttpUtils.handleRejectedFlow(flow);
        }
    }
}
