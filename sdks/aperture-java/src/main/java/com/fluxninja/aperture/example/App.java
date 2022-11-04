package com.fluxninja.aperture.example;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FeatureFlow;
import com.fluxninja.aperture.sdk.FlowStatus;
import io.grpc.ConnectivityState;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import spark.Spark;

import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

public class App {
    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";

    final private ApertureSDK apertureSDK;
    final private ManagedChannel channel;

    public App(ApertureSDK apertureSDK, ManagedChannel channel){
        this.apertureSDK = apertureSDK;
        this.channel = channel;
    }

    public static void main(String[] args) {
        String agentHost = System.getenv("FN_AGENT_HOST");
        if (agentHost == null) {
            agentHost = DEFAULT_AGENT_HOST;
        }
        String agentPort = System.getenv("FN_AGENT_PORT");
        if (agentPort == null) {
            agentPort = DEFAULT_AGENT_PORT;
        }

        String target = String.format("%s:%s", agentHost, agentPort);
        final ManagedChannel channel = ManagedChannelBuilder.forTarget(target).build();

        ApertureSDK apertureSDK;
        try {
            apertureSDK = ApertureSDK.builder()
                    .setHost(agentHost)
                    .setPort(Integer.parseInt(agentPort))
                    .setDuration(Duration.ofMillis(1000))
                    .build();
        } catch (ApertureSDKException e) {
            e.printStackTrace();
            return;
        }

        App app = new App(apertureSDK, channel);
        String appPort = System.getenv("FN_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }
        Spark.port(Integer.parseInt(appPort));
        Spark.get("/super", app::handleSuperAPI);
        Spark.get("/connected", app::handleConnectedAPI);
        Spark.get("/health", app::handleHealthAPI);
    }

    private String handleSuperAPI(spark.Request req, spark.Response res) {
        Map<String, String> labels = new HashMap<>();

        // do some business logic to collect labels
        labels.put("user", "kenobi");

        // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow.
        FeatureFlow flow = this.apertureSDK.startFlow("awesomeFeature", labels);

        // See whether flow was accepted by Aperture Agent.
        if (flow.accepted()) {
            // Simulate work being done
            try {
                res.status(202);
                Thread.sleep(2000);
                // Need to call end() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
                // The first argument captures whether the feature captured by the Flow was successful or resulted in an error.
                // The second argument is error message for further diagnosis.
                flow.end(FlowStatus.OK);
            } catch (InterruptedException | ApertureSDKException e) {
                e.printStackTrace();
            }
        } else {
            // Flow has been rejected by Aperture Agent.
            try {
                res.status(403);
                flow.end(FlowStatus.Error);
            } catch (ApertureSDKException e) {
                e.printStackTrace();
            }
        }
        return "";
    }

    private String handleConnectedAPI(spark.Request req, spark.Response res) {
        ConnectivityState state = this.channel.getState(true);
//        if (state.toString() != "READY") {
//            res.status(503);
//        }
        return state.toString();
    }

    private String handleHealthAPI(spark.Request req, spark.Response res) {
        return "Healthy";
    }
}
