package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpResponse;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.server.HttpService;
import com.linecorp.armeria.server.ServiceRequestContext;
import com.linecorp.armeria.server.SimpleDecoratingHttpService;

import java.util.Map;
import java.util.function.Function;

/**
 * Decorates an {@link HttpService} to enable flow control using provided {@link ApertureSDK}
 */
public class ApertureHTTPService extends SimpleDecoratingHttpService {
    private final ApertureSDK apertureSDK;

    public static Function<? super HttpService, ApertureHTTPService> newDecorator(ApertureSDK apertureSDK) {
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
        Map<String, String> labels = HttpUtils.labelsFromRequest(req);
        Flow flow = this.apertureSDK.startFlow("", labels);

        if (flow.accepted()) {
            HttpResponse res;
            try {
                res = unwrap().serve(ctx, req);
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
            HttpStatus code = HttpUtils.handleRejectedFlow(flow);
            return HttpResponse.of(code);
        }
    }
}
