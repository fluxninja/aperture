using Aperture.Flowcontrol.Check.V1;
using Grpc.Core;
using log4net;
using OpenTelemetry;
using OpenTelemetry.Trace;

namespace ApertureSDK.Core;

public sealed class ApertureSdkBuilder
{
    private readonly ILog _logger = LogManager.GetLogger(typeof(ApertureSdkBuilder));
    private string? _address;
    private string _agentApiKey = "";
    private ChannelCredentials _channelCredentials = new SslCredentials();
    private string _otlpExporterHeaders = "";
    private bool _useHttpsInOtlpExporter = true;


    public ApertureSdkBuilder SetAddress(string address)
    {
        _address = address;
        return this;
    }

    public ApertureSdkBuilder SetAgentApiKey(string agentApiKey)
    {
        _agentApiKey = agentApiKey;
        return this;
    }

    public ApertureSdkBuilder UseInsecureGrpc()
    {
        _channelCredentials = ChannelCredentials.Insecure;
        _useHttpsInOtlpExporter = false;
        return this;
    }

    public ApertureSdkBuilder SetChannelCredentials(ChannelCredentials channelCredentials)
    {
        _channelCredentials = channelCredentials;
        return this;
    }

    public ApertureSdkBuilder SetOtlpExporterHeaders(string headers)
    {
        _otlpExporterHeaders = headers;
        return this;
    }

    public IApertureSdk Build()
    {
        if (string.IsNullOrEmpty(_address))
        {
            _logger.Warn(
                $"Address not set when building Aperture SDK, defaulting to {Constants.DEFAULT_AGENT_ADDRESS}");
            _address = Constants.DEFAULT_AGENT_ADDRESS;
        }

        var otlpSpanExporterProtocol = _useHttpsInOtlpExporter ? "https" : "http";

        var channel = new Channel(_address, _channelCredentials);
        var flowControlClient = new FlowControlService.FlowControlServiceClient(channel);

        var tracerProvider = Sdk.CreateTracerProviderBuilder().AddOtlpExporter(
            opt =>
            {
                opt.Endpoint = new Uri($"{otlpSpanExporterProtocol}://{_address}");
                opt.Headers = _otlpExporterHeaders;
            }
        ).Build();

        var tracer = tracerProvider!.GetTracer(Constants.LIBRARY_NAME);

        return new ApertureSdk(flowControlClient, tracer, _agentApiKey);
    }
}
