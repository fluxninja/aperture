namespace ApertureSDK.Core;

public class FeatureFlowParams
{
    public FeatureFlowParams(
        string controlPoint,
        Dictionary<string, string> explicitLabels,
        bool rampMode,
        TimeSpan flowTimeout)
    {
        ControlPoint = controlPoint ?? throw new ArgumentNullException(nameof(controlPoint));
        ExplicitLabels = new Dictionary<string, string>(
            explicitLabels ?? throw new ArgumentNullException(nameof(explicitLabels)));
        RampMode = rampMode;
        FlowTimeout = flowTimeout;
    }

    public string ControlPoint { get; set; }
    public Dictionary<string, string> ExplicitLabels { get; set; }
    public TimeSpan FlowTimeout { get; set; }
    public bool RampMode { get; set; }
}
