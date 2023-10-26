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
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** A builder for configuring an {@link ApertureSDK}. */
public final class ApertureSDKBuilder {
    private String address;
    private boolean useHttpsInOtlpExporter = false;
    private boolean insecureGrpc = true;
    private String certFile;
    private final List<String> ignoredPaths;
    private boolean ignoredPathsMatchRegex = false;
    private TlsChannelCredentials.Builder tlsChannelCredentialsBuilder;
    private byte[] caCertFileContents;

    private static final Logger logger = LoggerFactory.getLogger(ApertureSDKBuilder.class);

    ApertureSDKBuilder() {
        ignoredPaths = new ArrayList<>();
        tlsChannelCredentialsBuilder = TlsChannelCredentials.newBuilder();
        caCertFileContents = null;
    }

    /**
     * Set address of Aperture Agent to connect to.
     *
     * @param address of Aperture Agent to connect to.
     * @return the builder object.
     */
    public ApertureSDKBuilder setAddress(String address) {
        this.address = address;
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
     * @throws IOException if custom root CA certificate file cannot be read.
     */
    public ApertureSDKBuilder setRootCertificateFile(String filename) throws IOException {
        this.certFile = filename;
        if (filename != null && !filename.isEmpty()) {
            this.tlsChannelCredentialsBuilder.trustManager(new File(this.certFile));
            this.caCertFileContents =
                    ByteStreams.toByteArray(Files.newInputStream(Paths.get(this.certFile)));
        }
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
    public ApertureSDK build() {
        String address = this.address;
        if (address == null) {
            logger.warn(
                    "Address not set when building Aperture SDK, defaulting to "
                            + DEFAULT_AGENT_ADDRESS);
            address = DEFAULT_AGENT_ADDRESS;
        }

        String OtlpSpanExporterProtocol = "http";
        if (this.useHttpsInOtlpExporter) {
            OtlpSpanExporterProtocol = "https";
        }

        ChannelCredentials creds;
        byte[] caCertContents = null;
        if (this.insecureGrpc) {
            creds = InsecureChannelCredentials.create();
        } else {
            if (this.certFile == null || this.certFile.isEmpty()) {
                creds = TlsChannelCredentials.create();
            } else {
                creds = this.tlsChannelCredentialsBuilder.build();
                caCertContents = this.caCertFileContents;
            }
        }

        ManagedChannel channel = Grpc.newChannelBuilder(address, creds).build();

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
                        .setEndpoint(String.format("%s://%s", OtlpSpanExporterProtocol, address))
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
                ignoredPaths,
                ignoredPathsMatchRegex);
    }
}
