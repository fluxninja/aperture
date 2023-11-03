package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowDecision;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RpcRequest;
import com.linecorp.armeria.common.RpcResponse;
import com.linecorp.armeria.server.RpcService;
import com.linecorp.armeria.server.ServiceRequestContext;
import com.linecorp.armeria.server.SimpleDecoratingRpcService;
import java.time.Duration;
import java.util.Map;
import java.util.function.Function;

/** Decorates an {@link RpcService} to enable flow control using provided {@link ApertureSDK} */
public class ApertureRPCService extends SimpleDecoratingRpcService {
    private final ApertureSDK apertureSDK;
    private final String controlPointName;
    private final boolean rampMode;
    private final Duration flowTimeout;

    public static Function<? super RpcService, ApertureRPCService> newDecorator(
            ApertureSDK apertureSDK, String controlPointName) {
        ApertureRPCServiceBuilder builder = new ApertureRPCServiceBuilder();
        builder.setApertureSDK(apertureSDK).setControlPointName(controlPointName);
        return builder::build;
    }

    public static Function<? super RpcService, ApertureRPCService> newDecorator(
            ApertureSDK apertureSDK, String controlPointName, boolean rampMode) {
        ApertureRPCServiceBuilder builder = new ApertureRPCServiceBuilder();
        builder.setApertureSDK(apertureSDK)
                .setControlPointName(controlPointName)
                .setEnableRampMode(rampMode);
        return builder::build;
    }

    public ApertureRPCService(
            RpcService delegate,
            ApertureSDK apertureSDK,
            String controlPointName,
            boolean rampMode,
            Duration flowTimeout) {
        super(delegate);
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
        this.rampMode = rampMode;
        this.flowTimeout = flowTimeout;
    }

    @Override
    public RpcResponse serve(ServiceRequestContext ctx, RpcRequest req) throws Exception {
        Map<String, String> labels = RpcUtils.labelsFromRequest(req);
        Flow flow =
                this.apertureSDK.startFlow(
                        this.controlPointName, labels, this.rampMode, flowTimeout);

        FlowDecision flowDecision = flow.getDecision();
        boolean flowAccepted =
                (flowDecision == FlowDecision.Accepted
                        || (flowDecision == FlowDecision.Unreachable && !this.rampMode));

        if (flowAccepted) {
            RpcResponse res;
            try {
                res = unwrap().serve(ctx, req);
            } catch (Exception e) {
                flow.setStatus(FlowStatus.Error);
                throw e;
            } finally {
                flow.end();
            }
            return res;
        } else {
            HttpStatus code = RpcUtils.handleRejectedFlow(flow);
            return RpcResponse.ofFailure(new Exception(code.toString()));
        }
    }
}
