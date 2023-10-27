package com.fluxninja.aperture.instrumentation;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKBuilder;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

public class Config {
    public static final String CONFIG_FILENAME_PROPERTY = "aperture.javaagent.config.file";

    public static final String AGENT_ADDRESS_PROPERTY = "aperture.agent.address";
    public static final String AGENT_API_KEY_PROPERTY = "aperture.agent.api.key";
    public static final String RAMP_MODE_PROPERTY = "aperture.javaagent.enable.ramp.mode";
    public static final String CONNECTION_TIMEOUT_MILLIS_PROPERTY =
            "aperture.connection.timeout.millis";
    public static final String CONTROL_POINT_NAME_PROPERTY = "aperture.control.point.name";
    public static final String IGNORED_PATHS_PROPERTY = "aperture.javaagent.ignored.paths";
    public static final String IGNORED_PATHS_REGEX_PROPERTY =
            "aperture.javaagent.ignored.paths.regex";
    public static final String INSECURE_GRPC_PROPERTY = "aperture.javaagent.insecure.grpc";
    public static final String ROOT_CERTIFICATE_FILE_PROPERTY =
            "aperture.javaagent.root.certificate";

    private static final String AGENT_ADDRESS_DEFAULT_VALUE = "localhost:8089";
    private static final String AGENT_API_KEY_DEFAULT_VALUE = "";
    private static final String RAMP_MODE_PROPERTY_DEFAULT_VALUE = "false";
    private static final String CONNECTION_TIMEOUT_MILLIS_DEFAULT_VALUE = "1000";
    private static final String IGNORED_PATHS_DEFAULT_VALUE = "";
    private static final String IGNORED_PATHS_REGEX_DEFAULT_VALUE = "false";
    private static final String INSECURE_GRPC_DEFAULT_VALUE = "true";
    private static final String ROOT_CERTIFICATE_FILE_DEFAULT_VALUE = "";

    private static final List<String> allProperties =
            new ArrayList<String>() {
                {
                    add(AGENT_ADDRESS_PROPERTY);
                    add(AGENT_API_KEY_PROPERTY);
                    add(RAMP_MODE_PROPERTY);
                    add(CONNECTION_TIMEOUT_MILLIS_PROPERTY);
                    add(CONTROL_POINT_NAME_PROPERTY);
                    add(IGNORED_PATHS_PROPERTY);
                    add(IGNORED_PATHS_REGEX_PROPERTY);
                    add(INSECURE_GRPC_PROPERTY);
                    add(ROOT_CERTIFICATE_FILE_PROPERTY);
                }
            };

    static Properties loadProperties() {
        Properties props = new Properties();
        String configFileName = System.getProperty(CONFIG_FILENAME_PROPERTY);
        if (configFileName == null) {
            configFileName = System.getenv(envNameFromPropertyName(CONFIG_FILENAME_PROPERTY));
        }
        try {
            if (configFileName != null) {
                props.load(Files.newInputStream(Paths.get(configFileName)));
            }
        } catch (IOException e) {
            throw new RuntimeException("Could not read properties from file", e);
        }

        // Get property overrides from env and commandline
        for (String key : allProperties) {
            String val = getFromEnv(key);
            if (val != null) {
                props.put(key, val);
            }
        }

        return props;
    }

    static String getFromEnv(String name) {
        // Read system property; If not set, use env variable.
        String systemProperty = System.getProperty(name);
        if (systemProperty != null) {
            return systemProperty;
        }
        String envVariableName = envNameFromPropertyName(name);
        return System.getenv(envVariableName);
    }

    public static ApertureSDKWrapper newSDKWrapperFromConfig() {
        ApertureSDKBuilder builder = ApertureSDK.builder();
        Properties config = loadProperties();
        ApertureSDK sdk;
        String controlPointName;
        boolean rampMode;
        Duration flowTimeout;
        try {
            controlPointName = config.getProperty(CONTROL_POINT_NAME_PROPERTY);
            rampMode =
                    Boolean.parseBoolean(
                            config.getProperty(
                                    RAMP_MODE_PROPERTY, RAMP_MODE_PROPERTY_DEFAULT_VALUE));

            flowTimeout =
                    Duration.ofMillis(
                            Integer.parseInt(
                                    config.getProperty(
                                            CONNECTION_TIMEOUT_MILLIS_PROPERTY,
                                            CONNECTION_TIMEOUT_MILLIS_DEFAULT_VALUE)));
            ApertureSDKBuilder sdkBuilder =
                    builder.setAddress(
                                    config.getProperty(
                                            AGENT_ADDRESS_PROPERTY, AGENT_ADDRESS_DEFAULT_VALUE))
                            .setAgentAPIKey(
                                    config.getProperty(
                                            AGENT_API_KEY_PROPERTY, AGENT_API_KEY_DEFAULT_VALUE))
                            .addIgnoredPaths(
                                    config.getProperty(
                                            IGNORED_PATHS_PROPERTY, IGNORED_PATHS_DEFAULT_VALUE))
                            .setIgnoredPathsMatchRegex(
                                    Boolean.parseBoolean(
                                            config.getProperty(
                                                    IGNORED_PATHS_REGEX_PROPERTY,
                                                    IGNORED_PATHS_REGEX_DEFAULT_VALUE)));

            boolean insecureGrpc =
                    Boolean.parseBoolean(
                            config.getProperty(
                                    INSECURE_GRPC_PROPERTY, INSECURE_GRPC_DEFAULT_VALUE));
            String caCertificateFile =
                    config.getProperty(
                            ROOT_CERTIFICATE_FILE_PROPERTY, ROOT_CERTIFICATE_FILE_DEFAULT_VALUE);

            sdkBuilder.useInsecureGrpc(insecureGrpc);
            if (!caCertificateFile.isEmpty()) {
                sdkBuilder.setRootCertificateFile(caCertificateFile);
            }

            sdk = sdkBuilder.build();
        } catch (Exception e) {
            throw new RuntimeException("failed to create Aperture SDK from config", e);
        }

        if (controlPointName == null || controlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must be set");
        }

        return new ApertureSDKWrapper(sdk, controlPointName, rampMode, flowTimeout);
    }

    private static String envNameFromPropertyName(String propertyName) {
        return propertyName.toUpperCase().replace(".", "_");
    }
}
