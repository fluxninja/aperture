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
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_SSL_CERT = "";

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
        String insecureGrpcString = System.getenv("FN_INSECURE_GRPC");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        System.setProperty("FN_INSECURE_GRPC", insecureGrpcString);

        String sslCertFile = System.getenv("FN_SSL_CERTIFICATE_FILE");
        if (sslCertFile == null) {
            sslCertFile = DEFAULT_SSL_CERT;
        }
        System.setProperty("FN_SSL_CERTIFICATE_FILE", sslCertFile);

        SpringApplication.run(SpringBootApp.class, args);
    }
}
