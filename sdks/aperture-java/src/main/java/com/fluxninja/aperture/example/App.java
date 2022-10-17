package com.fluxninja.aperture.example;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
import spark.Spark;

import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

public class App {
    public static final String DEFAULT_APP_PORT = "18080";
    public static final String DEFAULT_AGENT_HOST = "aperture-agent.aperture-agent.svc.cluster.local";
    public static final String DEFAULT_AGENT_PORT = "8089";
    final private ApertureSDK apertureSDK;
    public App(
            ApertureSDK apertureSDK
    ){
        this.apertureSDK = apertureSDK;
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

        App app = new App(apertureSDK);
        String appPort = System.getenv("FN_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }
        Spark.port(Integer.parseInt(appPort));
        Spark.get("/super", app::handleSuperAPI);
    }

    private String handleSuperAPI(spark.Request req, spark.Response res) {
        Map<String, String> labels = new HashMap<>();

        // do some business logic to collect labels
        labels.put("user", "kenobi");

        // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow.
        Flow flow = this.apertureSDK.startFlow("awesomeFeature", labels);

        // See whether flow was accepted by Aperture Agent.
        if (flow.accepted()) {
            // Simulate work being done
            try {
                Thread.sleep(5000);
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
                flow.end(FlowStatus.Error);
            } catch (ApertureSDKException e) {
                e.printStackTrace();
            }
        }
        return "Hello world";
    }
}
