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

    private static String getEnvOrDefault(String envVar, String defaultValue) {
        String value = System.getenv(envVar);
        if (value == null) {
            value = defaultValue;
        }
        System.setProperty(envVar, value);
        return value;
    }

    public static void main(String[] args) {
        getEnvOrDefault("APERTURE_AGENT_ADDRESS", DEFAULT_AGENT_ADDRESS);
        getEnvOrDefault("APERTURE_AGENT_API_KEY", "");
        getEnvOrDefault("APERTURE_CONTROL_POINT_NAME", DEFAULT_CONTROL_POINT_NAME);
        getEnvOrDefault("APERTURE_ENABLE_RAMP_MODE", DEFAULT_RAMP_MODE);
        getEnvOrDefault("APERTURE_GRPC_TIMEOUT_MS", DEFAULT_GRPC_TIMEOUT_MS);
        getEnvOrDefault("APERTURE_AGENT_INSECURE", DEFAULT_INSECURE_GRPC);
        getEnvOrDefault("APERTURE_ROOT_CERTIFICATE_FILE", DEFAULT_ROOT_CERT);

        SpringApplication.run(SpringBootApp.class, args);
    }
}
