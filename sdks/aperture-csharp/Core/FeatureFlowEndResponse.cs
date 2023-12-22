namespace ApertureSDK.Core;

using Aperture.Flowcontrol.Check.V1;

public class FeatureFlowEndResponse
{
    public FeatureFlowEndResponse(
        Exception? error,
        FlowEndResponse? flowEndResponse)
    {
        Error = error;
        FlowEndResponse = flowEndResponse;
    }

    public Exception? Error { get; set; }
    public FlowEndResponse? FlowEndResponse { get; set; }
}
