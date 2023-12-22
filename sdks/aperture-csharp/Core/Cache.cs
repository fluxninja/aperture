namespace ApertureSDK.Core;

using Aperture.Flowcontrol.Check.V1;

public enum LookupStatus
{
    HIT,
    MISS
}

public static class CacheLookupStatusConverter
{
    public static LookupStatus ConvertCacheLookupStatus(CacheLookupStatus status)
    {
        return status == CacheLookupStatus.Hit ? LookupStatus.HIT : LookupStatus.MISS;
    }
}

public static class CacheErrorConverter
{
    public static Exception? ConvertCacheError(string errorMessage)
    {
        return string.IsNullOrEmpty(errorMessage) ? null : new Exception(errorMessage);
    }
}

public class LookupResponse
{
    public object? Value { get; }
    public LookupStatus LookupStatus { get; }
    public Exception? Error { get; }

    public LookupResponse(object? value, LookupStatus lookupStatus, Exception? error)
    {
        Value = value;
        LookupStatus = lookupStatus;
        Error = error;
    }
}

public class UpsertResponse
{
    public Exception? Error { get; }

    public UpsertResponse(Exception? error)
    {
        Error = error;
    }
}

public class DeleteResponse
{
    public Exception? Error { get; }

    public DeleteResponse(Exception? error)
    {
        Error = error;
    }
}
