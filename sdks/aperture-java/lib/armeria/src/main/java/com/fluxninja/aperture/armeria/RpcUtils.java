package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RpcRequest;
import java.util.HashMap;
import java.util.Map;

class RpcUtils {
    protected static HttpStatus handleRejectedFlow(Flow flow) {
        try {
            flow.end(FlowStatus.Unset);
        } catch (ApertureSDKException e) {
            e.printStackTrace();
        }
        if (flow.checkResponse() != null) {
            CheckResponse.RejectReason reason = flow.checkResponse().getRejectReason();
            switch (reason) {
                case REJECT_REASON_RATE_LIMITED:
                    return HttpStatus.TOO_MANY_REQUESTS;
                case REJECT_REASON_NO_TOKENS:
                    return HttpStatus.SERVICE_UNAVAILABLE;
                case REJECT_REASON_REGULATED:
                    return HttpStatus.FORBIDDEN;
                case REJECT_REASON_NONE:
                case UNRECOGNIZED:
            }
        }
        return HttpStatus.FORBIDDEN;
    }

    protected static Map<String, String> labelsFromRequest(RpcRequest req) {
        Map<String, String> labels = new HashMap<>();
        labels.put("rpc.method", req.method());
        return labels;
    }
}
