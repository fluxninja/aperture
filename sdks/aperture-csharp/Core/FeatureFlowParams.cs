namespace ApertureSDK.Core;

using Google.Protobuf.Collections;
using Grpc.Core;

public class FeatureFlowParams
{
    public FeatureFlowParams(
        string controlPoint,
        Dictionary<string, string> explicitLabels,
        bool rampMode,
        TimeSpan flowTimeout,
        CallOptions callOptions,
        string resultCacheKey,
        RepeatedField<string> globalCacheKeys)
    {
        ControlPoint = controlPoint ?? throw new ArgumentNullException(nameof(controlPoint));
        ExplicitLabels = new Dictionary<string, string>(
            explicitLabels ?? throw new ArgumentNullException(nameof(explicitLabels)));
        RampMode = rampMode;
        FlowTimeout = flowTimeout;
        CallOptions = callOptions;
        ResultCacheKey = resultCacheKey;
        GlobalCacheKeys = globalCacheKeys;
    }

    public string ControlPoint { get; set; }
    public Dictionary<string, string> ExplicitLabels { get; set; }
    public TimeSpan FlowTimeout { get; set; }
    public bool RampMode { get; set; }
    public CallOptions CallOptions { get; set; }
    public string ResultCacheKey { get; set; }
    public RepeatedField<string> GlobalCacheKeys { get; set; }
}
