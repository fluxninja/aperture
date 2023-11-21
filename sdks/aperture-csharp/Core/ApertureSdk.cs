using System.Text;
using System.Web;
using Aperture.Flowcontrol.Check.V1;
using Grpc.Core;
using log4net;
using OpenTelemetry;
using OpenTelemetry.Trace;

namespace ApertureSDK.Core;

public class ApertureSdk : IApertureSdk
{
    private readonly string? _apiKey;
    private readonly FlowControlService.FlowControlServiceClient _flowControlClient;
    private readonly ILog _logger = LogManager.GetLogger(typeof(ApertureSdk));
    private readonly Tracer _tracer;

    public ApertureSdk(
        FlowControlService.FlowControlServiceClient flowControlClient,
        Tracer tracer,
        string apiKey)
    {
        _flowControlClient = flowControlClient;
        _tracer = tracer;
        _apiKey = apiKey;
    }

    public IFlow StartFlow(FeatureFlowParams parameters)
    {
        var controlPoint = parameters.ControlPoint;
        var explicitLabels = parameters.ExplicitLabels;
        var rampMode = parameters.RampMode;
        var flowTimeout = parameters.FlowTimeout;
        var labels = new Dictionary<string, string>();

        foreach (var entry in Baggage.Current.GetBaggage())
        {
            string value;
            try
            {
                value = HttpUtility.UrlDecode(entry.Value, Encoding.UTF8);
            }
            catch (Exception e)
            {
                // This should never happen, as Encoding.UTF8 is a valid encoding
                _logger.Debug("URL decoding failed: {e}", e);
                throw;
            }

            labels[entry.Key] = value;
        }

        foreach (var explicitLabel in explicitLabels) labels[explicitLabel.Key] = explicitLabel.Value;

        var checkReq = new CheckRequest();
        checkReq.ControlPoint = controlPoint;
        checkReq.RampMode = rampMode;
        foreach (var label in labels) checkReq.Labels.Add(label.Key, label.Value);

        using var span = _tracer.StartSpan("Aperture Check");

        span.SetAttribute(Constants.FLOW_START_TIMESTAMP_LABEL, Utils.GetCurrentEpochNanos());
        span.SetAttribute(Constants.SOURCE_LABEL, "sdk");

        CheckResponse? res = null;
        try
        {
            var opts = new CallOptions();
            if (flowTimeout != TimeSpan.Zero) opts = opts.WithDeadline(DateTime.UtcNow.Add(flowTimeout));
            if (_apiKey != null)
            {
                var headers = new Metadata();
                headers.Add("x-api-key", _apiKey);
                opts = opts.WithHeaders(headers);
            }

            res = _flowControlClient.Check(checkReq, opts);
        }
        catch (Exception e)
        {
            // Deadline exceeded or couldn't reach the agent - request should not be blocked
            _logger.Debug("Check call caused an exception");
        }

        span.SetAttribute(Constants.WORKLOAD_START_TIMESTAMP_LABEL, Utils.GetCurrentEpochNanos());

        return new FeatureFlow(res, span, false, rampMode);
    }

    public static ApertureSdkBuilder Builder()
    {
        return new ApertureSdkBuilder();
    }
}
