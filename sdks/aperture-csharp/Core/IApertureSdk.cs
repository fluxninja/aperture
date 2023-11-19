namespace ApertureSDK.Core;

/// <summary>
///     Defines the interface for the Aperture SDK.
/// </summary>
public interface IApertureSdk
{
    /// <summary>
    ///     Starts a new feature flow with the given parameters.
    /// </summary>
    /// <param name="parameters">The parameters for the feature flow.</param>
    /// <returns>An IFlow instance representing the started feature flow.</returns>
    IFlow StartFlow(FeatureFlowParams parameters);
}
