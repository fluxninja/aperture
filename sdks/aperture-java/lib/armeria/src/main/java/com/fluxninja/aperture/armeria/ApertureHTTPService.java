package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpResponse;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.server.HttpService;
import com.linecorp.armeria.server.ServiceRequestContext;
import com.linecorp.armeria.server.SimpleDecoratingHttpService;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Function;

/** Decorates an {@link HttpService} to enable flow control using provided {@link ApertureSDK} */
public class ApertureHTTPService extends SimpleDecoratingHttpService {
    private final ApertureSDK apertureSDK;

    public static Function<? super HttpService, ApertureHTTPService> newDecorator(
            ApertureSDK apertureSDK) {
        ApertureHTTPServiceBuilder builder = new ApertureHTTPServiceBuilder();
        builder.setApertureSDK(apertureSDK);
        return builder::build;
    }

    public ApertureHTTPService(HttpService delegate, ApertureSDK apertureSDK) {
        super(delegate);
        this.apertureSDK = apertureSDK;
    }

    @Override
    public HttpResponse serve(ServiceRequestContext ctx, HttpRequest req) throws Exception {
        AttributeContext attributes = HttpUtils.attributesFromRequest(req);
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(req.path(), attributes);

        if (flow.ignored()) {
            return unwrap().serve(ctx, req);
        }

        if (flow.accepted()) {
            HttpResponse res;
            try {
                List<HeaderValueOption> newHeaders = new ArrayList<>();
                if (flow.checkResponse() != null) {
                    newHeaders = flow.checkResponse().getOkResponse().getHeadersList();
                }
                HttpRequest newRequest = HttpUtils.updateHeaders(req, newHeaders);
                ctx.updateRequest(newRequest);

                res = unwrap().serve(ctx, newRequest);
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                // ending flow failed
                e.printStackTrace();
                return HttpResponse.of(HttpStatus.INTERNAL_SERVER_ERROR);
            } catch (Exception e) {
                try {
                    flow.end(FlowStatus.Error);
                } catch (ApertureSDKException ae) {
                    e.printStackTrace();
                    ae.printStackTrace();
                }
                throw e;
            }
            return res;
        } else {
            HttpStatus code = HttpUtils.handleRejectedFlow(flow);
            return HttpResponse.of(code);
        }
    }
}
