package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RpcRequest;
import java.util.HashMap;
import java.util.Map;

class RpcUtils {
    protected static HttpStatus handleRejectedFlow(Flow flow) {
        flow.setStatus(FlowStatus.Unset);
        flow.end();
        return HttpStatus.valueOf(flow.getRejectionHttpStatusCode());
    }

    protected static Map<String, String> labelsFromRequest(RpcRequest req) {
        Map<String, String> labels = new HashMap<>();
        labels.put("rpc.method", req.method());
        return labels;
    }
}
