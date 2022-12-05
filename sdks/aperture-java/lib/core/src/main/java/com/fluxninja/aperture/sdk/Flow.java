package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.google.protobuf.util.JsonFormat;
import io.opentelemetry.api.trace.Span;

import static com.fluxninja.aperture.sdk.Constants.*;

public final class Flow {
  private final CheckResponse checkResponse;
  private final Span span;
  private boolean ended;

  Flow(
      CheckResponse checkResponse,
      Span span,
      boolean ended) {
    this.checkResponse = checkResponse;
    this.span = span;
    this.ended = ended;
  }

  public boolean accepted() {
    if (this.checkResponse == null) {
      return true;
    }
    return this.checkResponse.getDecisionType() == CheckResponse.DecisionType.DECISION_TYPE_ACCEPTED;
  }

  public CheckResponse checkResponse() {
    return this.checkResponse;
  }

  public void end(FlowStatus statusCode) throws ApertureSDKException {
    if (this.ended) {
      throw new ApertureSDKException("Flow already ended");
    }
    this.ended = true;

    String checkResponseJSONBytes;
    try {
      checkResponseJSONBytes = JsonFormat.printer().print(this.checkResponse);
    } catch (com.google.protobuf.InvalidProtocolBufferException e) {
      throw new ApertureSDKException(e);
    }

    this.span.setAttribute(FLOW_STATUS_LABEL, statusCode.name())
        .setAttribute(CHECK_RESPONSE_LABEL, checkResponseJSONBytes)
        .setAttribute(FLOW_STOP_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

    this.span.end();
  }
}
