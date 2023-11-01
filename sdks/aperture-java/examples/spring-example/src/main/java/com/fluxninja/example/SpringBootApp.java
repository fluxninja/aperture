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
        String controlPointName = System.getenv("APERTURE_CONTROL_POINT_NAME");
        if (controlPointName == null) {
            controlPointName = DEFAULT_CONTROL_POINT_NAME;
        }
        System.setProperty("APERTURE_CONTROL_POINT_NAME", controlPointName);
        String rampMode = System.getenv("APERTURE_ENABLE_RAMP_MODE");
        if (rampMode == null) {
            rampMode = DEFAULT_RAMP_MODE;
        }
        System.setProperty("APERTURE_ENABLE_RAMP_MODE", rampMode);
        String grpcTimeoutMs = System.getenv("APERTURE_GRPC_TIMEOUT_MS");
        if (grpcTimeoutMs == null) {
            grpcTimeoutMs = DEFAULT_GRPC_TIMEOUT_MS;
        }
        System.setProperty("APERTURE_GRPC_TIMEOUT_MS", grpcTimeoutMs);
        String insecureGrpcString = System.getenv("APERTURE_AGENT_INSECURE");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        System.setProperty("APERTURE_AGENT_INSECURE", insecureGrpcString);

        String rootCertFile = System.getenv("APERTURE_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }
        System.setProperty("APERTURE_ROOT_CERTIFICATE_FILE", rootCertFile);

        SpringApplication.run(SpringBootApp.class, args);
    }
}
