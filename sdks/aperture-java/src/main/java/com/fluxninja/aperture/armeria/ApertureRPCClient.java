package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.linecorp.armeria.client.ClientRequestContext;
import com.linecorp.armeria.client.RpcClient;
import com.linecorp.armeria.client.SimpleDecoratingRpcClient;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RpcRequest;
import com.linecorp.armeria.common.RpcResponse;

import java.util.Map;
import java.util.function.Function;

/**
 * Decorates an {@link RpcClient} to enable flow control using provided
 * {@link ApertureSDK}
 */
public class ApertureRPCClient extends SimpleDecoratingRpcClient {
  private final ApertureSDK apertureSDK;

  public static Function<? super RpcClient, ApertureRPCClient> newDecorator(ApertureSDK apertureSDK) {
    ApertureRPCClientBuilder builder = new ApertureRPCClientBuilder();
    builder.setApertureSDK(apertureSDK);
    return builder::build;
  }

  public ApertureRPCClient(RpcClient delegate, ApertureSDK apertureSDK) {
    super(delegate);
    this.apertureSDK = apertureSDK;
  }

  @Override
  public RpcResponse execute(ClientRequestContext ctx, RpcRequest req) throws Exception {
    Map<String, String> labels = RpcUtils.labelsFromRequest(req);
    Flow flow = this.apertureSDK.startFlow("", labels);

    if (flow.accepted()) {
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
