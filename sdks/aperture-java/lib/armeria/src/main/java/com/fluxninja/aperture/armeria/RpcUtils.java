package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.EndResponse;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RpcRequest;
import java.util.HashMap;
import java.util.Map;

class RpcUtils {
    protected static HttpStatus handleRejectedFlow(Flow flow) {
        flow.setStatus(FlowStatus.Unset);
        EndResponse endResponse = flow.end();
        if (endResponse.getError() != null) {
            System.err.println("Error ending flow: " + endResponse.getError().getMessage());
        }

        return HttpStatus.valueOf(flow.getRejectionHttpStatusCode());
    }

    protected static Map<String, String> labelsFromRequest(RpcRequest req) {
        Map<String, String> labels = new HashMap<>();
        labels.put("rpc.method", req.method());
        return labels;
    }
}
