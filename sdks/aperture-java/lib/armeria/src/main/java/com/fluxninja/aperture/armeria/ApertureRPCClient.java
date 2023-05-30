package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.*;
import com.linecorp.armeria.client.ClientRequestContext;
import com.linecorp.armeria.client.RpcClient;
import com.linecorp.armeria.client.SimpleDecoratingRpcClient;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RpcRequest;
import com.linecorp.armeria.common.RpcResponse;
import java.util.Map;
import java.util.function.Function;

/** Decorates an {@link RpcClient} to enable flow control using provided {@link ApertureSDK} */
public class ApertureRPCClient extends SimpleDecoratingRpcClient {
    private final ApertureSDK apertureSDK;
    private final String controlPointName;
    private final boolean failOpen;

    public static Function<? super RpcClient, ApertureRPCClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName) {
        ApertureRPCClientBuilder builder = new ApertureRPCClientBuilder();
        builder.setApertureSDK(apertureSDK).setControlPointName(controlPointName);
        return builder::build;
    }

    public static Function<? super RpcClient, ApertureRPCClient> newDecorator(
            ApertureSDK apertureSDK, String controlPointName, boolean failOpen) {
        ApertureRPCClientBuilder builder = new ApertureRPCClientBuilder();
        builder.setApertureSDK(apertureSDK)
                .setControlPointName(controlPointName)
                .setEnableFailOpen(failOpen);
        return builder::build;
    }

    public ApertureRPCClient(
            RpcClient delegate,
            ApertureSDK apertureSDK,
            String controlPointName,
            boolean failOpen) {
        super(delegate);
        this.apertureSDK = apertureSDK;
        this.controlPointName = controlPointName;
        this.failOpen = failOpen;
    }

    @Override
    public RpcResponse execute(ClientRequestContext ctx, RpcRequest req) throws Exception {
        Map<String, String> labels = RpcUtils.labelsFromRequest(req);
        Flow flow = this.apertureSDK.startFlow(this.controlPointName, labels);

        FlowDecision flowDecision = flow.getDecision();
        boolean flowAccepted =
                (flowDecision == FlowDecision.Accepted
                        || (flowDecision == FlowDecision.Unreachable && this.failOpen));

        if (flowAccepted) {
            RpcResponse res;
            try {
                res = unwrap().execute(ctx, req);
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                // ending flow failed
                e.printStackTrace();
                return RpcResponse.ofFailure(e);
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
            HttpStatus code = RpcUtils.handleRejectedFlow(flow);
            return RpcResponse.ofFailure(new Exception(code.toString()));
        }
    }
}
