package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.FlowControlServiceGrpc;
import com.fluxninja.generated.envoy.service.auth.v3.AuthorizationGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.SimpleSpanProcessor;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static com.fluxninja.aperture.sdk.Constants.DEFAULT_RPC_TIMEOUT;
import static com.fluxninja.aperture.sdk.Constants.LIBRARY_NAME;

/** A builder for configuring an {@link ApertureSDK}. */
public final class ApertureSDKBuilder {
  private Duration timeout;
  private String host;
  private int port;
  private boolean useHttps = false;
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

  public ApertureSDKBuilder useHttps() {
    this.useHttps = true;
    return this;
  }

  /** Adds comma-separated paths to ignore in traffic control points.
   * @param paths comma-separated list of paths to ignore when creating traffic control points.
   * @return the builder object.
   */
  public ApertureSDKBuilder addBlockedPaths(String paths) {
    if (paths == null || paths.isEmpty()) {
      return this;
    }
    return this.addBlockedPaths(Arrays.asList(paths.split("\\s*,\\s*")));
  }

  /** Adds paths to ignore in traffic control points.
   * @param paths list of paths to ignore when creating traffic control points.
   * @return the builder object.
   */
  public ApertureSDKBuilder addBlockedPaths(List<String> paths) {
    this.blockedPaths.addAll(paths);
    return this;
  }

  /** Whether paths should be matched by regex. If false, exact matches will be expected.
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

    String protocol = "http";
    if (this.useHttps) {
      protocol = "https";
    }

    Duration timeout = this.timeout;
    if (timeout == null) {
      timeout = DEFAULT_RPC_TIMEOUT;
    }

    ManagedChannel channel = ManagedChannelBuilder.forAddress(host, port).usePlaintext().build();
    FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient = FlowControlServiceGrpc
        .newBlockingStub(channel);
    AuthorizationGrpc.AuthorizationBlockingStub envoyAuthzClient = AuthorizationGrpc.newBlockingStub(channel);

    OtlpGrpcSpanExporter spanExporter = OtlpGrpcSpanExporter.builder()
        .setEndpoint(String.format("%s://%s:%d", protocol, host, port))
        .build();
    SdkTracerProvider traceProvider = SdkTracerProvider.builder()
        .addSpanProcessor(SimpleSpanProcessor.create(spanExporter))
        .build();
    Tracer tracer = traceProvider.tracerBuilder(LIBRARY_NAME).build();

    return new ApertureSDK(
        flowControlClient,
        envoyAuthzClient,
        tracer,
        timeout,
        blockedPaths,
        blockedPathsMatchRegex);
  }
}
