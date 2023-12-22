using System.Net;
using System.Text;
using ApertureSDK.Core;
using Google.Protobuf.Collections;
using Google.Protobuf.WellKnownTypes;
using log4net;
using log4net.Config;

[assembly: XmlConfigurator(ConfigFile = "log4net.config", Watch = true)]
XmlConfigurator.Configure();

var log = LogManager.GetLogger("AppLogger");

var agentAddress = Environment.GetEnvironmentVariable("APERTURE_AGENT_ADDRESS") ?? "localhost:8089";
var apiKey = Environment.GetEnvironmentVariable("APERTURE_API_KEY") ?? "";
var insecureGrpcString = Environment.GetEnvironmentVariable("APERTURE_AGENT_INSECURE") ?? "true";
var featureName = Environment.GetEnvironmentVariable("APERTURE_FEATURE_NAME") ?? "awesome-feature";
var appPort = Environment.GetEnvironmentVariable("APERTURE_APP_PORT") ?? "8080";

var builder = ApertureSdk
    .Builder()
    .SetAddress(agentAddress)
    .SetApiKey(apiKey);

if (bool.Parse(insecureGrpcString))
    builder
        .UseInsecureGrpc();
var sdk = builder.Build();

var url = $"http://*:{appPort}/";
log.Info($"Starting server at {url}");

using (var listener = new HttpListener())
{
    listener.Prefixes.Add(url);
    listener.Start();

    while (true)
    {
        var context = listener.GetContext();
        var request = context.Request;
        var response = context.Response;

        if (request.Url!.AbsolutePath == "/super")
        {
            var labels = new Dictionary<string, string>();
            labels.Add("user", "kenobi");
            var pms = new FeatureFlowParams(
                featureName,
                labels,
                false,
                TimeSpan.FromSeconds(5),
                new Grpc.Core.CallOptions(),
                "test",
                new RepeatedField<string> { "test" });
            var flow = sdk.StartFlow(pms);
            if (flow.ShouldRun())
            {
                Thread.Sleep(2000);
                SimpleHandlePath((int)HttpStatusCode.OK, "Hello world!", response);
            }
            else
            {
                SimpleHandlePath(flow.GetRejectionHttpStatusCode(), "REJECTED!", response);
            }

            flow.End();
        }
        else if (request.Url.AbsolutePath == "/super2")
        {
            // START: handleRequest
            // do some business logic to collect labels
            var labels = new Dictionary<string, string>();
            labels.Add("userId", "some_user_id");
            labels.Add("userTier", "gold");
            labels.Add("priority", "100");

            var rampMode = false;
            var flowTimeout = TimeSpan.FromSeconds(5);
            var pms = new FeatureFlowParams(
                "featureName",
                labels,
                rampMode,
                flowTimeout,
                new Grpc.Core.CallOptions(),
                "test",
                new RepeatedField<string> { "test" });
            var flow = sdk.StartFlow(pms);
            if (flow.ShouldRun())
            {
                // do actual work
                Thread.Sleep(2000);
                SimpleHandlePath((int)HttpStatusCode.OK, "Hello world!", response);
            }
            else
            {
                // handle flow rejection by Aperture Agent
                flow.SetStatus(FlowStatus.Error);
                SimpleHandlePath(flow.GetRejectionHttpStatusCode(), "REJECTED!", response);
            }

            var endResponse = flow.End();
            if (endResponse.Error != null)
            {
                // handle end failure
                log.Error("Failed to end flow: {e}", endResponse.Error);
            }
            else if (endResponse.FlowEndResponse != null)
            {
                // handle end success
                log.Info("Ended flow with response: " + endResponse.FlowEndResponse.ToString());
            }
            // END: handleRequest
        }
        else if (request.Url.AbsolutePath == "/notsuper")
        {
            SimpleHandlePath((int)HttpStatusCode.OK, "Hello world!", response);
        }
        else if (request.Url.AbsolutePath == "/health")
        {
            SimpleHandlePath((int)HttpStatusCode.OK, "Healthy", response);
        }
        else if (request.Url.AbsolutePath == "/connected")
        {
            SimpleHandlePath((int)HttpStatusCode.OK, "Connected", response);
        }
        else
        {
            SimpleHandlePath((int)HttpStatusCode.NotFound, "Not found", response);
        }

        response.Close();
    }
}

void SimpleHandlePath(int statusCode, string responseString, HttpListenerResponse response)
{
    var buffer = Encoding.UTF8.GetBytes(responseString);

    response.StatusCode = statusCode;
    response.ContentLength64 = buffer.Length;
    response.OutputStream.Write(buffer, 0, buffer.Length);
}
