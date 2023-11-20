using NodaTime;

namespace ApertureSDK.Core;

public static class Utils
{
    public static long GetCurrentEpochNanos()
    {
        var now = SystemClock.Instance.GetCurrentInstant();
        return now.ToUnixTimeTicks() * 100;
    }
}
