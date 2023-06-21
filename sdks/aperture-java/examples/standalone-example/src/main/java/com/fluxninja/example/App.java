package com.fluxninja.example;

import com.fluxninja.aperture.sdk.*;
import io.grpc.ConnectivityState;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.io.IOException;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import spark.Spark;

public class App {
    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";
    public static final String DEFAULT_FEATURE_NAME = "awesome_feature";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

    public static final Logger logger = LoggerFactory.getLogger(App.class);

    private final ApertureSDK apertureSDK;
    private final ManagedChannel channel;
    private final String featureName;

    public App(ApertureSDK apertureSDK, ManagedChannel channel, String featureName) {
        this.apertureSDK = apertureSDK;
        this.channel = channel;
        this.featureName = featureName;
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
        String insecureGrpcString = System.getenv("FN_INSECURE_GRPC");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String rootCertFile = System.getenv("FN_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }

        String target = String.format("%s:%s", agentHost, agentPort);
        final ManagedChannel channel = ManagedChannelBuilder.forTarget(target).build();

        ApertureSDK apertureSDK;
        try {
            apertureSDK =
                    ApertureSDK.builder()
                            .setHost(agentHost)
                            .setPort(Integer.parseInt(agentPort))
                            .setFlowTimeout(Duration.ofMillis(1000))
                            .useInsecureGrpc(insecureGrpc)
                            .setRootCertificateFile(rootCertFile)
                            .build();
        } catch (IOException e) {
            e.printStackTrace();
            return;
        }

        String featureName = System.getenv("FN_FEATURE_NAME");
        if (featureName == null) {
            featureName = DEFAULT_FEATURE_NAME;
        }

        App app = new App(apertureSDK, channel, featureName);
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

        // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a
        // Flow.
        Flow flow = this.apertureSDK.startFlow(this.featureName, labels);

        // See whether flow was accepted by Aperture Agent.
        try {
            if (flow.shouldRun()) {
                // Simulate work being done
                res.status(202);
                Thread.sleep(2000);
            } else {
                // Flow has been rejected by Aperture Agent.
                res.status(flow.getRejectionHttpStatusCode());
            }
        } catch (Exception e) {
            // Flow Status captures whether the feature captured by the Flow was
            // successful or resulted in an error. When not explicitly set,
            // the default value is FlowStatus.OK .
            flow.setStatus(FlowStatus.Error);
            logger.error("Error in flow execution", e);
        } finally {
            flow.end();
        }
        return "";
    }

    private String handleConnectedAPI(spark.Request req, spark.Response res) {
        ConnectivityState state = this.channel.getState(true);
        // if (state.toString() != "READY") {
        // res.status(503);
        // }
        return state.toString();
    }

    private String handleHealthAPI(spark.Request req, spark.Response res) {
        return "Healthy";
    }
}
