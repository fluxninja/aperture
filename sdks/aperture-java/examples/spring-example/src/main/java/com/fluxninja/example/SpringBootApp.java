package com.fluxninja.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.servlet.ServletComponentScan;

@ServletComponentScan
@SpringBootApplication
public class SpringBootApp {
    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";
    public static final String DEFAULT_FAIL_OPEN = "true";
    public static final String DEFAULT_CONTROL_POINT_NAME = "awesome_feature";
    public static final String DEFAULT_GRPC_TIMEOUT_MS = "1000";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

    public static void main(String[] args) {
        String agentHost = System.getenv("FN_AGENT_HOST");
        if (agentHost == null) {
            agentHost = DEFAULT_AGENT_HOST;
        }
        System.setProperty("FN_AGENT_HOST", agentHost);
        String agentPort = System.getenv("FN_AGENT_PORT");
        if (agentPort == null) {
            agentPort = DEFAULT_AGENT_PORT;
        }
        System.setProperty("FN_AGENT_PORT", agentPort);
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
        String failOpen = System.getenv("FN_ENABLE_FAIL_OPEN");
        if (failOpen == null) {
            failOpen = DEFAULT_FAIL_OPEN;
        }
        System.setProperty("FN_ENABLE_FAIL_OPEN", failOpen);
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
