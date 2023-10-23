package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.*;
import com.linecorp.armeria.client.ClientRequestContext;
import com.linecorp.armeria.client.RpcClient;
import com.linecorp.armeria.client.SimpleDecoratingRpcClient;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RpcRequest;
import com.linecorp.armeria.common.RpcResponse;
import java.time.Duration;
import java.util.Map;
import java.util.function.Function;

/** Decorates an {@link RpcClient} to enable flow control using provided {@link ApertureSDK} */
public class ApertureRPCClient extends SimpleDecoratingRpcClient {
    private final ApertureSDK apertureSDK;
    private final String controlPointName;
    private final boolean rampMode;
    private final Duration flowTimeout;

    public static Function<? super RpcClient, ApertureRPCClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName) {
        ApertureRPCClientBuilder builder = new ApertureRPCClientBuilder();
        builder.setApertureSDK(apertureSDK).setControlPointName(controlPointName);
        return builder::build;
    }

    public static Function<? super RpcClient, ApertureRPCClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName, boolean rampMode) {
        ApertureRPCClientBuilder builder = new ApertureRPCClientBuilder();
        builder.setApertureSDK(apertureSDK)
                .setControlPointName(controlPointName)
                .setEnableRampMode(rampMode);
        return builder::build;
    }

    public ApertureRPCClient(
            RpcClient delegate,
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
    public RpcResponse execute(ClientRequestContext ctx, RpcRequest req) throws Exception {
        Map<String, String> labels = RpcUtils.labelsFromRequest(req);
        Flow flow = this.apertureSDK.startFlow(this.controlPointName, labels, false, flowTimeout);

        FlowDecision flowDecision = flow.getDecision();
        boolean flowAccepted =
                (flowDecision == FlowDecision.Accepted
                        || (flowDecision == FlowDecision.Unreachable && !this.rampMode));

        if (flowAccepted) {
            RpcResponse res;
            try {
                res = unwrap().execute(ctx, req);
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
