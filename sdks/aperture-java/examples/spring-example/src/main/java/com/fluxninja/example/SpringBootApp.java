package com.fluxninja.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.servlet.ServletComponentScan;

@ServletComponentScan
@SpringBootApplication
public class SpringBootApp {
    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_ADDRESS = "localhost:8089";
    public static final String DEFAULT_RAMP_MODE = "false";
    public static final String DEFAULT_CONTROL_POINT_NAME = "awesome_feature";
    public static final String DEFAULT_GRPC_TIMEOUT_MS = "1000";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

    public static void main(String[] args) {
        String agentAddress = System.getenv("APERTURE_AGENT_ADDRESS");
        if (agentAddress == null) {
            agentAddress = DEFAULT_AGENT_ADDRESS;
        }
        System.setProperty("APERTURE_AGENT_ADDRESS", agentAddress);
        String agentAPIKey = System.getenv("APERTURE_AGENT_API_KEY");
        if (agentAPIKey == null) {
            agentAPIKey = "";
        }
        System.setProperty("APERTURE_AGENT_API_KEY", agentAPIKey);
        String appPort = System.getenv("FN_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }
        System.setProperty("FN_APP_PORT", appPort);
        String controlPointName = System.getenv("FN_CONTROL_POINT_NAME");
        if (controlPointName == null) {
            controlPointName = DEFAULT_CONTROL_POINT_NAME;
        }
        System.setProperty("FN_CONTROL_POINT_NAME", controlPointName);
        String rampMode = System.getenv("FN_ENABLE_RAMP_MODE");
        if (rampMode == null) {
            rampMode = DEFAULT_RAMP_MODE;
        }
        System.setProperty("FN_ENABLE_RAMP_MODE", rampMode);
        String grpcTimeoutMs = System.getenv("FN_GRPC_TIMEOUT_MS");
        if (grpcTimeoutMs == null) {
            grpcTimeoutMs = DEFAULT_GRPC_TIMEOUT_MS;
        }
        System.setProperty("FN_GRPC_TIMEOUT_MS", grpcTimeoutMs);
        String insecureGrpcString = System.getenv("FN_INSECURE_GRPC");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        System.setProperty("FN_INSECURE_GRPC", insecureGrpcString);

        String rootCertFile = System.getenv("FN_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }
        System.setProperty("FN_ROOT_CERTIFICATE_FILE", rootCertFile);

        SpringApplication.run(SpringBootApp.class, args);
    }
}
