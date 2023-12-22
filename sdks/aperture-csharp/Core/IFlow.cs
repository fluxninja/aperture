namespace ApertureSDK.Core;

public interface IFlow
{
    bool ShouldRun();

    FeatureFlowEndResponse End();

    void SetStatus(FlowStatus status);

    int GetRejectionHttpStatusCode();
}
