using System.Net;
using Aperture.Flowcontrol.Check.V1;
using Google.Protobuf;
using Grpc.Core;
using log4net;
using OpenTelemetry.Trace;

namespace ApertureSDK.Core;

public class FeatureFlow : IFlow
{
    private readonly CheckResponse? _checkResponse;
    private readonly ILog _logger = LogManager.GetLogger(typeof(FeatureFlow));
    private readonly bool _rampMode;
    private readonly TelemetrySpan _span;
    private bool _ended;
    private FlowStatus _flowStatus;
    private FlowControlService.FlowControlServiceClient _flowControlServiceClient;
    private CallOptions _callOptions;

    public FeatureFlow(CheckResponse? checkResponse, TelemetrySpan span, bool ended, bool rampMode, FlowControlService.FlowControlServiceClient flowControlServiceClient, CallOptions callOptions)
    {
        _checkResponse = checkResponse;
        _span = span;
        _ended = ended;
        _rampMode = rampMode;
        _flowStatus = FlowStatus.Ok;
        _flowControlServiceClient = flowControlServiceClient;
        _callOptions = callOptions;
    }

    public bool ShouldRun()
    {
        return GetDecision() == FlowDecision.Accepted || (GetDecision() == FlowDecision.Unreachable && !_rampMode);
    }

    public FeatureFlowEndResponse End()
    {
        if (_ended)
        {
            _logger.Warn("Attempting to end an already ended Flow");
            return new FeatureFlowEndResponse(new Exception("Attempting to end an already ended Flow"), null);
        }

        _ended = true;

        var checkResponseJsonBytes = string.Empty;

        try
        {
            if (_checkResponse != null) checkResponseJsonBytes = JsonFormatter.Default.Format(_checkResponse);
        }
        catch (InvalidProtocolBufferException e)
        {
            _logger.Warn("Could not attach check response when ending flow: {e}", e);
        }

        _logger.Debug($"Ending a Flow with status {_flowStatus}");

        _span
            .SetAttribute(Constants.FLOW_STATUS_LABEL, _flowStatus.ToString())
            .SetAttribute(Constants.CHECK_RESPONSE_LABEL, checkResponseJsonBytes)
            .SetAttribute(Constants.FLOW_STOP_TIMESTAMP_LABEL, Utils.GetCurrentEpochNanos());

        _span.End();

        var inflightRequestRef = new List<InflightRequestRef>();

        for (var i = 0; i < _checkResponse?.LimiterDecisions.Count; i++)
        {
            var inflightRequest = new InflightRequestRef
            {
                PolicyHash = _checkResponse.LimiterDecisions[i].PolicyHash,
                PolicyName = _checkResponse.LimiterDecisions[i].PolicyName,
                ComponentId = _checkResponse.LimiterDecisions[i].ComponentId
            };

            if (_checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo != null) {
                inflightRequest.Label = _checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo.Label;

                if (_checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo.TokensInfo != null)
                {
                    inflightRequest.Tokens = _checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo.TokensInfo.Consumed;
                }
            }

            inflightRequestRef.Add(inflightRequest);
        }

        if (inflightRequestRef.Count > 0)
        {
            var flowEndRequest = new FlowEndRequest
            {
                ControlPoint = _checkResponse!.ControlPoint
            };
            flowEndRequest.InflightRequests.AddRange(inflightRequestRef);

            try
            {
                var flowEndResponse = _flowControlServiceClient.FlowEnd(flowEndRequest, _callOptions);
                return new FeatureFlowEndResponse(null, flowEndResponse);
            }
            catch (Exception e)
            {
                _logger.Warn("Could not end flow: {e}", e);
                return new FeatureFlowEndResponse(e, null);
            }

        }

        return new FeatureFlowEndResponse(null, null);
    }

    public void SetStatus(FlowStatus status)
    {
        if (_ended) _logger.Warn("Attempting to change status of an already ended Flow");

        _flowStatus = status;
    }

    public int GetRejectionHttpStatusCode()
    {
        if (GetDecision() == FlowDecision.Rejected)
            switch (_checkResponse!.RejectReason)
            {
                case CheckResponse.Types.RejectReason.RateLimited:
                    return (int)HttpStatusCode.TooManyRequests;
                case CheckResponse.Types.RejectReason.NoTokens:
                    return (int)HttpStatusCode.ServiceUnavailable;
                case CheckResponse.Types.RejectReason.NotSampled:
                    return (int)HttpStatusCode.Forbidden;
                case CheckResponse.Types.RejectReason.NoMatchingRamp:
                    return (int)HttpStatusCode.Forbidden;
                default:
                    return (int)HttpStatusCode.Forbidden;
            }

        throw new InvalidOperationException("Flow not rejected");
    }

    private FlowDecision GetDecision()
    {
        if (_checkResponse == null) return FlowDecision.Unreachable;

        if (_checkResponse.DecisionType == CheckResponse.Types.DecisionType.Accepted) return FlowDecision.Accepted;

        return FlowDecision.Rejected;
    }
}
