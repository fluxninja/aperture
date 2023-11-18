namespace ApertureSDK.Core;

public interface IFlow
{
    bool ShouldRun();

    void End();

    void SetStatus(FlowStatus status);

    int GetRejectionHttpStatusCode();
}
