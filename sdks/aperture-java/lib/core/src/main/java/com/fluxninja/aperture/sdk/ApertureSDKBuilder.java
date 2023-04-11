package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.DEFAULT_RPC_TIMEOUT;
import static com.fluxninja.aperture.sdk.Constants.LIBRARY_NAME;

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

/** A builder for configuring an {@link ApertureSDK}. */
public final class ApertureSDKBuilder {
    private Duration timeout;
    private String host;
    private int port;
    private boolean useHttpsInOtlpExporter = false;
    private boolean insecureGrpc = true;
    private String certFile;
    private final List<String> blockedPaths;
    private boolean blockedPathsMatchRegex = false;

    ApertureSDKBuilder() {
        blockedPaths = new ArrayList<>();
    }

    public ApertureSDKBuilder setHost(String host) {
        this.host = host;
        return this;
    }

    public ApertureSDKBuilder setPort(int port) {
        this.port = port;
        return this;
    }

    public ApertureSDKBuilder setDuration(Duration timeout) {
        this.timeout = timeout;
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
     * @param paths comma-separated list of paths to ignore when creating traffic control points.
     * @return the builder object.
     */
    public ApertureSDKBuilder addBlockedPaths(String paths) {
        if (paths == null || paths.isEmpty()) {
            return this;
        }
        return this.addBlockedPaths(Arrays.asList(paths.split("\\s*,\\s*")));
    }

    /**
     * Adds paths to ignore in traffic control points.
     *
     * @param paths list of paths to ignore when creating traffic control points.
     * @return the builder object.
     */
    public ApertureSDKBuilder addBlockedPaths(List<String> paths) {
        this.blockedPaths.addAll(paths);
        return this;
    }

    /**
     * Whether paths should be matched by regex. If false, exact matches will be expected.
     *
     * @param flag whether paths should be matched by regex.
     * @return the builder object.
     */
    public ApertureSDKBuilder setBlockedPathMatchRegex(boolean flag) {
        this.blockedPathsMatchRegex = flag;
        return this;
    }

    public ApertureSDK build() throws ApertureSDKException {
        String host = this.host;
        if (host == null) {
            throw new ApertureSDKException("host needs to be set");
        }

        int port = this.port;
        if (port == 0) {
            throw new ApertureSDKException("port needs to be set");
        }

        String OtlpSpanExporterProtocol = "http";
        if (this.useHttpsInOtlpExporter) {
            OtlpSpanExporterProtocol = "https";
        }

        Duration timeout = this.timeout;
        if (timeout == null) {
            timeout = DEFAULT_RPC_TIMEOUT;
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
                timeout,
                blockedPaths,
                blockedPathsMatchRegex);
    }
}
