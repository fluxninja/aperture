using System.Net;
using Aperture.Flowcontrol.Check.V1;
using Google.Protobuf;
using Grpc.Core;
using Google.Protobuf.Collections;
using log4net;
using OpenTelemetry.Trace;
using Google.Protobuf.WellKnownTypes;

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
    private readonly string _resultCacheKey;
    private readonly RepeatedField<string> _globalCacheKeys;

    public FeatureFlow(CheckResponse? checkResponse, TelemetrySpan span, bool ended, bool rampMode, FlowControlService.FlowControlServiceClient flowControlServiceClient, CallOptions callOptions, string resultCacheKey, RepeatedField<string> globalCacheKeys)
    {
        _checkResponse = checkResponse;
        _span = span;
        _ended = ended;
        _rampMode = rampMode;
        _flowStatus = FlowStatus.Ok;
        _flowControlServiceClient = flowControlServiceClient;
        _callOptions = callOptions;
        _resultCacheKey = resultCacheKey;
        _globalCacheKeys = globalCacheKeys;
    }

    public bool ShouldRun()
    {
        return GetDecision() == FlowDecision.Accepted || (GetDecision() == FlowDecision.Unreachable && !_rampMode);
    }

    public LookupResponse ResultCache()
    {
        if (_checkResponse == null)
        {
            return new LookupResponse(null, LookupStatus.MISS, new Exception("check response is null"));
        }

        if (!ShouldRun())
        {
            return new LookupResponse(null, LookupStatus.MISS, new Exception("flow was rejected"));
        }

        if (_checkResponse.CacheLookupResponse == null || _checkResponse.CacheLookupResponse.ResultCacheResponse == null)
        {
            return new LookupResponse(null, LookupStatus.MISS, new Exception("result cache is null"));
        }

        var cacheLookupResponse = _checkResponse.CacheLookupResponse.ResultCacheResponse;
        return new LookupResponse(
            cacheLookupResponse.Value.ToStringUtf8(),
            CacheLookupStatusConverter.ConvertCacheLookupStatus(cacheLookupResponse.LookupStatus),
            CacheErrorConverter.ConvertCacheError(cacheLookupResponse.Error));
    }

    public UpsertResponse SetResultCache(Object value, Duration ttl)
    {
        if (_resultCacheKey == "")
        {
            return new UpsertResponse(new Exception("result cache key is empty"));
        }

        if (_checkResponse == null)
        {
            return new UpsertResponse(new Exception("check response is null"));
        }

        var request = new CacheUpsertRequest
        {
            ControlPoint = _checkResponse.ControlPoint,
            ResultCacheEntry = new CacheEntry
            {
                Key = _resultCacheKey,
                Value = ByteString.CopyFromUtf8(value.ToString()),
                Ttl = ttl
            }
        };

        try
        {
            var cacheUpsertResponse = _flowControlServiceClient.CacheUpsert(request, _callOptions);
            return new UpsertResponse(null);
        }
        catch (Exception e)
        {
            _logger.Warn("Could not set result cache: {e}", e);
            return new UpsertResponse(e);
        }

    }

    public DeleteResponse DeleteResultCache()
    {
        if (_resultCacheKey == "")
        {
            return new DeleteResponse(new Exception("result cache key is empty"));
        }

        if (_checkResponse == null)
        {
            return new DeleteResponse(new Exception("check response is null"));
        }

        var request = new CacheDeleteRequest
        {
            ControlPoint = _checkResponse.ControlPoint,
            ResultCacheKey = _resultCacheKey
        };

        try
        {
            var cacheDeleteResponse = _flowControlServiceClient.CacheDelete(request, _callOptions);
            return new DeleteResponse(null);
        }
        catch (Exception e)
        {
            _logger.Warn("Could not delete result cache: {e}", e);
            return new DeleteResponse(e);
        }
    }

    public LookupResponse GlobalCache(string key)
    {
        if (_checkResponse == null)
        {
            return new LookupResponse(null, LookupStatus.MISS, new Exception("check response is null"));
        }

        if (!ShouldRun())
        {
            return new LookupResponse(null, LookupStatus.MISS, new Exception("flow was rejected"));
        }

        if (_checkResponse.CacheLookupResponse == null || _checkResponse.CacheLookupResponse.GlobalCacheResponses == null)
        {
            return new LookupResponse(null, LookupStatus.MISS, new Exception("global cache is null"));
        }

        var cacheLookupResponses = _checkResponse.CacheLookupResponse.GlobalCacheResponses;
        var cacheLookupResponse = cacheLookupResponses.GetValueOrDefault(key);

        if (cacheLookupResponse == null)
        {
            return new LookupResponse(null, LookupStatus.MISS, new Exception("global cache is null"));
        }

        return new LookupResponse(
            cacheLookupResponse.Value.ToStringUtf8(),
            CacheLookupStatusConverter.ConvertCacheLookupStatus(cacheLookupResponse.LookupStatus),
            CacheErrorConverter.ConvertCacheError(cacheLookupResponse.Error));
    }

    public UpsertResponse SetGlobalCache(string key, Object value, Duration ttl)
    {
        if (_checkResponse == null)
        {
            return new UpsertResponse(new Exception("check response is null"));
        }

        var request = new CacheUpsertRequest
        {
            ControlPoint = _checkResponse.ControlPoint,
            GlobalCacheEntries = { }
        };
        request.GlobalCacheEntries.Add(key, new CacheEntry
        {
            Value = ByteString.CopyFromUtf8(value.ToString()),
            Ttl = ttl
        });

        try
        {
            var cacheUpsertResponse = _flowControlServiceClient.CacheUpsert(request, _callOptions);
            return new UpsertResponse(null);
        }
        catch (Exception e)
        {
            _logger.Warn("Could not set global cache: {e}", e);
            return new UpsertResponse(e);
        }
    }

    public DeleteResponse DeleteGlobalCache(string key)
    {
        if (_checkResponse == null)
        {
            return new DeleteResponse(new Exception("check response is null"));
        }

        var request = new CacheDeleteRequest
        {
            ControlPoint = _checkResponse.ControlPoint,
            GlobalCacheKeys = { key }
        };

        try
        {
            var cacheDeleteResponse = _flowControlServiceClient.CacheDelete(request, _callOptions);
            return new DeleteResponse(null);
        }
        catch (Exception e)
        {
            _logger.Warn("Could not delete global cache: {e}", e);
            return new DeleteResponse(e);
        }
    }

    public FeatureFlowEndResponse End()
    {
        if (_ended)
        {
            _logger.Warn("Attempting to end an already ended Flow");
            return new FeatureFlowEndResponse(new Exception("Attempting to end an already ended Flow"), null);
        }

        if (_checkResponse == null)
        {
            _logger.Warn("Attempting to end a Flow without a check response");
            return new FeatureFlowEndResponse(new Exception("Attempting to end a Flow without a check response"), null);
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

            if (_checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo != null)
            {
                var inflightRequest = new InflightRequestRef
                {
                    PolicyHash = _checkResponse.LimiterDecisions[i].PolicyHash,
                    PolicyName = _checkResponse.LimiterDecisions[i].PolicyName,
                    ComponentId = _checkResponse.LimiterDecisions[i].ComponentId,
                    Label = _checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo.Label,
                    RequestId = _checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo.RequestId
                };

                if (_checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo.TokensInfo != null)
                {
                    inflightRequest.Tokens = _checkResponse.LimiterDecisions[i].ConcurrencyLimiterInfo.TokensInfo.Consumed;
                }
                inflightRequestRef.Add(inflightRequest);
            }

            if (_checkResponse.LimiterDecisions[i].ConcurrencySchedulerInfo != null)
            {
                var inflightRequest = new InflightRequestRef
                {
                    PolicyHash = _checkResponse.LimiterDecisions[i].PolicyHash,
                    PolicyName = _checkResponse.LimiterDecisions[i].PolicyName,
                    ComponentId = _checkResponse.LimiterDecisions[i].ComponentId,
                    Label = _checkResponse.LimiterDecisions[i].ConcurrencySchedulerInfo.Label,
                    RequestId = _checkResponse.LimiterDecisions[i].ConcurrencySchedulerInfo.RequestId
                };

                if (_checkResponse.LimiterDecisions[i].ConcurrencySchedulerInfo.TokensInfo != null)
                {
                    inflightRequest.Tokens = _checkResponse.LimiterDecisions[i].ConcurrencySchedulerInfo.TokensInfo.Consumed;
                }
                inflightRequestRef.Add(inflightRequest);
            }

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
