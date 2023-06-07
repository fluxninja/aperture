package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.*;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.FlowControlServiceGrpc;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.FlowControlServiceHTTPGrpc;
import com.google.common.io.ByteStreams;
import io.grpc.*;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter;
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporterBuilder;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.SimpleSpanProcessor;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** A builder for configuring an {@link ApertureSDK}. */
public final class ApertureSDKBuilder {
    private Duration flowTimeout;
    private String host;
    private int port;
    private boolean useHttpsInOtlpExporter = false;
    private boolean insecureGrpc = true;
    private String certFile;
    private final List<String> ignoredPaths;
    private boolean ignoredPathsMatchRegex = false;

    private static Logger logger = LoggerFactory.getLogger(ApertureSDKBuilder.class);

    ApertureSDKBuilder() {
        ignoredPaths = new ArrayList<>();
    }

    /**
     * Set hostname of Aperture Agent to connect to.
     *
     * @param host hostname of Aperture Agent to connect to.
     * @return the builder object.
     */
    public ApertureSDKBuilder setHost(String host) {
        this.host = host;
        return this;
    }

    /**
     * Set port number of Aperture Agent to connect to.
     *
     * @param port port number of Aperture Agent to connect to.
     * @return the builder object.
     */
    public ApertureSDKBuilder setPort(int port) {
        this.port = port;
        return this;
    }

    /**
     * Set timeout for connection to Aperture Agent. Set to 0 to block until response is received.
     *
     * @param timeout timeout for connection to Aperture Agent.
     * @return the builder object.
     */
    public ApertureSDKBuilder setFlowTimeout(Duration timeout) {
        this.flowTimeout = timeout;
        return this;
    }

    /**
     * Makes the span exporter use https instead of http.
     *
     * @deprecated Use {@link #useInsecureGrpc(boolean)} instead.
     * @return the builder object.
     */
    @Deprecated
    public ApertureSDKBuilder useHttps() {
        this.useHttpsInOtlpExporter = true;
        return this;
    }

    /**
     * Sets custom root CA certificate to be used by SSL connection.
     *
     * @param filename path to file containing custom root CA certificate.
     * @return the builder object.
     */
    public ApertureSDKBuilder setRootCertificateFile(String filename) {
        this.certFile = filename;
        return this;
    }

    /**
     * Makes the SDK use plaintext if true, and SSL/TLS if false. Custom root CA certificates can be
     * passes using {@link #setRootCertificateFile(String)} method.
     *
     * @return the builder object.
     */
    public ApertureSDKBuilder useInsecureGrpc(boolean enabled) {
        this.insecureGrpc = enabled;
        this.useHttpsInOtlpExporter = !enabled;
        return this;
    }

    /**
     * Adds comma-separated paths to ignore in traffic control points.
     *
     * @deprecated use {@link #addIgnoredPaths(String)}
     * @param paths comma-separated list of paths to ignore when creating traffic control points.
     * @return the builder object.
     */
    public ApertureSDKBuilder addBlockedPaths(String paths) {
        if (paths == null || paths.isEmpty()) {
            return this;
        }
        return this.addIgnoredPaths(Arrays.asList(paths.split("\\s*,\\s*")));
    }

    /**
     * Adds comma-separated paths to ignore in traffic control points.
     *
     * @param paths comma-separated list of paths to ignore when creating traffic control points.
     * @return the builder object.
     */
    public ApertureSDKBuilder addIgnoredPaths(String paths) {
        if (paths == null || paths.isEmpty()) {
            return this;
        }
        return this.addIgnoredPaths(Arrays.asList(paths.split("\\s*,\\s*")));
    }

    /**
     * Adds paths to ignore in traffic control points.
     *
     * @deprecated use {@link #addIgnoredPaths(List)}
     * @param paths list of paths to ignore when creating traffic control points.
     * @return the builder object.
     */
    public ApertureSDKBuilder addBlockedPaths(List<String> paths) {
        this.ignoredPaths.addAll(paths);
        return this;
    }

    /**
     * Adds paths to ignore in traffic control points.
     *
     * @param paths list of paths to ignore when creating traffic control points.
     * @return the builder object.
     */
    public ApertureSDKBuilder addIgnoredPaths(List<String> paths) {
        this.ignoredPaths.addAll(paths);
        return this;
    }

    /**
     * Whether ignored paths should be matched by regex. If false, exact matches will be expected.
     *
     * @deprecated use {@link #setIgnoredPathsMatchRegex(boolean)}
     * @param flag whether paths should be matched by regex.
     * @return the builder object.
     */
    public ApertureSDKBuilder setBlockedPathMatchRegex(boolean flag) {
        this.ignoredPathsMatchRegex = flag;
        return this;
    }

    /**
     * Whether ignored paths should be matched by regex. If false, exact matches will be expected.
     *
     * @param flag whether paths should be matched by regex.
     * @return the builder object.
     */
    public ApertureSDKBuilder setIgnoredPathsMatchRegex(boolean flag) {
        this.ignoredPathsMatchRegex = flag;
        return this;
    }

    /**
     * Build an ApertureSDK object using the configured parameters.
     *
     * @return The constructed ApertureSDK object.
     */
    public ApertureSDK build() throws ApertureSDKException {
        String host = this.host;
        if (host == null) {
            logger.warn(
                    "Host not set when building Aperture SDK, defaulting to " + DEFAULT_AGENT_HOST);
            host = DEFAULT_AGENT_HOST;
        }

        int port = this.port;
        if (port == 0) {
            logger.warn(
                    "Port not set when building Aperture SDK, defaulting to " + DEFAULT_AGENT_PORT);
            port = DEFAULT_AGENT_PORT;
        }

        String OtlpSpanExporterProtocol = "http";
        if (this.useHttpsInOtlpExporter) {
            OtlpSpanExporterProtocol = "https";
        }

        Duration flowTimeout = this.flowTimeout;
        if (flowTimeout == null) {
            flowTimeout = DEFAULT_RPC_TIMEOUT;
        }

        ChannelCredentials creds;
        byte[] caCertContents = null;
        if (this.insecureGrpc) {
            creds = InsecureChannelCredentials.create();
        } else {
            if (this.certFile == null || this.certFile.isEmpty()) {
                creds = TlsChannelCredentials.create();
            } else {
                try {
                    creds =
                            TlsChannelCredentials.newBuilder()
                                    .trustManager(new File(this.certFile))
                                    .build();
                    caCertContents =
                            ByteStreams.toByteArray(Files.newInputStream(Paths.get(this.certFile)));
                } catch (IOException e) {
                    // cert file not found
                    throw new ApertureSDKException(e);
                }
            }
        }

        String target = host + ":" + port;
        ManagedChannel channel = Grpc.newChannelBuilder(target, creds).build();

        FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient =
                FlowControlServiceGrpc.newBlockingStub(channel);
        FlowControlServiceHTTPGrpc.FlowControlServiceHTTPBlockingStub httpFlowControlClient =
                FlowControlServiceHTTPGrpc.newBlockingStub(channel);

        OtlpGrpcSpanExporterBuilder spanExporterBuilder = OtlpGrpcSpanExporter.builder();
        if (caCertContents != null) {
            spanExporterBuilder.setTrustedCertificates(caCertContents);
        }

        OtlpGrpcSpanExporter spanExporter =
                spanExporterBuilder
                        .setEndpoint(
                                String.format("%s://%s:%d", OtlpSpanExporterProtocol, host, port))
                        .build();
        SdkTracerProvider traceProvider =
                SdkTracerProvider.builder()
                        .addSpanProcessor(SimpleSpanProcessor.create(spanExporter))
                        .build();
        Tracer tracer = traceProvider.tracerBuilder(LIBRARY_NAME).build();

        return new ApertureSDK(
                flowControlClient,
                httpFlowControlClient,
                tracer,
                flowTimeout,
                ignoredPaths,
                ignoredPathsMatchRegex);
    }
}
