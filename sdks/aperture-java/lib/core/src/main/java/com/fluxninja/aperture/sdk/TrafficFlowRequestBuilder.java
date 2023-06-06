package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import java.util.Map;

/**
 * The TrafficFlowRequestBuilder class is responsible for constructing TrafficFlowRequest objects.
 * It provides methods to set various parameters of the request such as control point, source,
 * destination, HTTP method, HTTP path, HTTP host, HTTP scheme, HTTP size, HTTP protocol, and HTTP
 * headers. Once all the desired parameters are set, the build() method can be called to create a
 * TrafficFlowRequest object.
 */
public class TrafficFlowRequestBuilder {
    private final CheckHTTPRequest.Builder checkHTTPRequestBuilder;
    private final CheckHTTPRequest.HttpRequest.Builder httpRequestBuilder;

    /**
     * Constructs a new TrafficFlowRequestBuilder object. Initializes the internal
     * CheckHTTPRequest.Builder and CheckHTTPRequest.HttpRequest.Builder objects.
     */
    TrafficFlowRequestBuilder() {
        checkHTTPRequestBuilder = CheckHTTPRequest.newBuilder();
        httpRequestBuilder = CheckHTTPRequest.HttpRequest.newBuilder();
    }

    /**
     * Sets the name of the control point of the traffic flow.
     *
     * @param controlPoint The control point name to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setControlPoint(String controlPoint) {
        checkHTTPRequestBuilder.setControlPoint(controlPoint);
        return this;
    }

    /**
     * Sets the source address of the traffic flow request.
     *
     * @param host The host name or IP address of the source.
     * @param port The port number of the source.
     * @param protocol The protocol used for the source.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setSource(String host, int port, String protocol) {
        checkHTTPRequestBuilder.setSource(Utils.createSocketAddress(host, port, protocol));
        return this;
    }

    /**
     * Sets the destination address of the traffic flow request.
     *
     * @param host The host name or IP address of the destination.
     * @param port The port number of the destination.
     * @param protocol The protocol used for the destination.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setDestination(String host, int port, String protocol) {
        checkHTTPRequestBuilder.setDestination(Utils.createSocketAddress(host, port, protocol));
        return this;
    }

    /**
     * Sets the HTTP method of the traffic flow request.
     *
     * @param httpMethod The HTTP method to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setHttpMethod(String httpMethod) {
        httpRequestBuilder.setMethod(httpMethod);
        return this;
    }

    /**
     * Sets the HTTP path of the traffic flow request.
     *
     * @param httpPath The HTTP path to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setHttpPath(String httpPath) {
        httpRequestBuilder.setPath(httpPath);
        return this;
    }

    /**
     * Sets the HTTP host of the traffic flow request.
     *
     * @param httpHost The HTTP host to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setHttpHost(String httpHost) {
        httpRequestBuilder.setHost(httpHost);
        return this;
    }

    /**
     * Sets the HTTP scheme of the traffic flow request.
     *
     * @param httpScheme The HTTP scheme to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setHttpScheme(String httpScheme) {
        httpRequestBuilder.setScheme(httpScheme);
        return this;
    }

    /**
     * Sets the HTTP size of the traffic flow request.
     *
     * @param httpSize The HTTP size to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setHttpSize(long httpSize) {
        httpRequestBuilder.setSize(httpSize);
        return this;
    }

    /**
     * Sets the HTTP protocol of the traffic flow request.
     *
     * @param httpProtocol The HTTP protocol to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setHttpProtocol(String httpProtocol) {
        httpRequestBuilder.setProtocol(httpProtocol);
        return this;
    }

    /**
     * Sets the HTTP headers of the traffic flow request.
     *
     * @param httpHeaders The map of HTTP headers to be set.
     * @return The TrafficFlowRequestBuilder object itself.
     */
    public TrafficFlowRequestBuilder setHttpHeaders(Map<String, String> httpHeaders) {
        httpRequestBuilder.putAllHeaders(httpHeaders);
        return this;
    }

    /**
     * Builds a TrafficFlowRequest object using the configured parameters.
     *
     * @return The constructed TrafficFlowRequest object.
     */
    public TrafficFlowRequest build() {
        checkHTTPRequestBuilder.setRequest(httpRequestBuilder.build());
        CheckHTTPRequest checkHTTPRequest = checkHTTPRequestBuilder.build();
        return new TrafficFlowRequest(checkHTTPRequest);
    }
}
