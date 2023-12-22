using Google.Protobuf.WellKnownTypes;

namespace ApertureSDK.Core;

public interface IFlow
{
    bool ShouldRun();

    FeatureFlowEndResponse End();

    void SetStatus(FlowStatus status);

    int GetRejectionHttpStatusCode();

    LookupResponse ResultCache();

    UpsertResponse SetResultCache(object value, Duration ttl);

    DeleteResponse DeleteResultCache();

    LookupResponse GlobalCache(string key);

    UpsertResponse SetGlobalCache(string key, object value, Duration ttl);

    DeleteResponse DeleteGlobalCache(string key);
}
