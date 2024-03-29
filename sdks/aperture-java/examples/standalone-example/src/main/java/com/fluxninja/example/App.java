package com.fluxninja.example;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.EndResponse;
import com.fluxninja.aperture.sdk.FeatureFlowParameters;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
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
    public static final String DEFAULT_AGENT_ADDRESS = "localhost:8089";
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
        String agentAddress = System.getenv("APERTURE_AGENT_ADDRESS");
        if (agentAddress == null) {
            agentAddress = DEFAULT_AGENT_ADDRESS;
        }
        String apiKey = System.getenv("APERTURE_API_KEY");
        if (apiKey == null) {
            apiKey = "";
        }
        String insecureGrpcString = System.getenv("APERTURE_AGENT_INSECURE");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String rootCertFile = System.getenv("APERTURE_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }

        final ManagedChannel channel = ManagedChannelBuilder.forTarget(agentAddress).build();

        // START: StandaloneExampleSDKInit

        ApertureSDK apertureSDK;
        try {
            apertureSDK =
                    ApertureSDK.builder()
                            .setAddress(agentAddress)
                            .setAPIKey(apiKey)
                            .useInsecureGrpc(insecureGrpc)
                            .setRootCertificateFile(rootCertFile)
                            .build();
        } catch (IOException e) {
            e.printStackTrace();
            return;
        }

        // END: StandaloneExampleSDKInit

        String featureName = System.getenv("APERTURE_FEATURE_NAME");
        if (featureName == null) {
            featureName = DEFAULT_FEATURE_NAME;
        }

        App app = new App(apertureSDK, channel, featureName);
        String appPort = System.getenv("APERTURE_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }
        Spark.port(Integer.parseInt(appPort));
        Spark.get("/super", app::handleSuperAPI);
        Spark.get("/super2", app::handleSuper2API);
        Spark.get("/connected", app::handleConnectedAPI);
        Spark.get("/health", app::handleHealthAPI);
    }

    private String handleSuperAPI(spark.Request req, spark.Response res) {
        Map<String, String> labels = new HashMap<>();

        // do some business logic to collect labels
        labels.put("user", "kenobi");

        FeatureFlowParameters params =
                FeatureFlowParameters.newBuilder("awesomeFeature")
                        .setExplicitLabels(labels)
                        .setRampMode(false)
                        .setFlowTimeout(Duration.ofMillis(1000))
                        .build();
        // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a
        // Flow.
        Flow flow = this.apertureSDK.startFlow(params);

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
            EndResponse endResponse = flow.end();
            if (endResponse.getError() != null) {
                logger.error("Error ending flow", endResponse.getError());
            }

            logger.info("Flow End response: {}", endResponse.getFlowEndResponse());
        }
        return "";
    }

    private String handleSuper2API(spark.Request req, spark.Response res) {

        // START: StandaloneExampleFlow

        Map<String, String> labels = new HashMap<>();

        // business logic produces labels
        labels.put("userId", "some_user_id");
        labels.put("userTier", "gold");
        labels.put("priority", "100");

        Boolean rampMode = false;

        FeatureFlowParameters params =
                FeatureFlowParameters.newBuilder("featureName")
                        .setExplicitLabels(labels)
                        .setRampMode(rampMode)
                        .setFlowTimeout(Duration.ofMillis(1000))
                        .build();
        // StartFlow performs a flowcontrolv1.Check call to Aperture. It returns a Flow.
        Flow flow = this.apertureSDK.startFlow(params);

        // See whether flow was accepted by Aperture.
        try {
            if (flow.shouldRun()) {
                // do actual work
                res.status(202);
            } else {
                // handle flow rejection by Aperture
                res.status(flow.getRejectionHttpStatusCode());
            }
        } catch (Exception e) {
            // Flow Status captures whether the feature captured by the Flow was
            // successful or resulted in an error. When not explicitly set,
            // the default value is FlowStatus.OK .
            flow.setStatus(FlowStatus.Error);
            logger.error("Error in flow execution", e);
        } finally {
            EndResponse endResponse = flow.end();
            if (endResponse.getError() != null) {
                logger.error("Error ending flow", endResponse.getError());
            }

            logger.info("Flow End response: {}", endResponse.getFlowEndResponse());
        }

        // END: StandaloneExampleFlow

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
