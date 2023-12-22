package com.fluxninja.aperture.sdk;

import java.time.Duration;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class FeatureFlowParameters {
    private String controlPoint;
    private Map<String, String> explicitLabels;
    private Boolean rampMode;
    private Duration flowTimeout;
    private String resultCacheKey;

    private List<String> globalCacheKeys;

    public static Builder newBuilder(String controlPoint) {
        return new Builder(controlPoint);
    }

    private FeatureFlowParameters() {
        // private constructor to enforce the use of the builder
    }

    public String getControlPoint() {
        return controlPoint;
    }

    public Map<String, String> getExplicitLabels() {
        return explicitLabels;
    }

    public Boolean getRampMode() {
        return rampMode;
    }

    public Duration getFlowTimeout() {
        return flowTimeout;
    }

    public String getResultCacheKey() {
        return resultCacheKey;
    }

    public List<String> getGlobalCacheKeys() {
        return globalCacheKeys;
    }

    public static class Builder {
        private final FeatureFlowParameters params;

        // Constructor to initialize the required parameter
        public Builder(String controlPoint) {
            params = new FeatureFlowParameters();
            params.controlPoint = controlPoint;
            params.explicitLabels = new HashMap<>();
            params.rampMode = false;
            params.flowTimeout = Constants.DEFAULT_RPC_TIMEOUT;
            params.resultCacheKey = "";
            params.globalCacheKeys = Collections.emptyList();
        }

        /**
         * Set the explicit labels for the FeatureFlowParameters.
         *
         * @param explicitLabels Labels sent to Aperture Agent
         * @return This builder for method chaining
         */
        public Builder setExplicitLabels(Map<String, String> explicitLabels) {
            params.explicitLabels = new HashMap<>(explicitLabels);
            return this;
        }

        /**
         * Set whether the flow should require a ramp component match.
         *
         * @param rampMode Whether the flow should require ramp component match
         * @return This builder for method chaining
         */
        public Builder setRampMode(Boolean rampMode) {
            params.rampMode = rampMode;
            return this;
        }

        /**
         * Set the timeout for the connection to Aperture Agent.
         *
         * @param flowTimeout Timeout for connection to Aperture Agent. Set to 0 to block until
         *     response
         * @return This builder for method chaining
         */
        public Builder setFlowTimeout(Duration flowTimeout) {
            params.flowTimeout = flowTimeout;
            return this;
        }

        /**
         * Set the result cache key for result cache request
         *
         * @param resultCacheKey The result cache key
         * @return This builder for method chaining
         */
        public Builder setResultCacheKey(String resultCacheKey) {
            params.resultCacheKey = resultCacheKey;
            return this;
        }

        /**
         * Add global cache keys for global cache request
         *
         * @param globalCacheKeys The global cache keys to add
         * @return This builder for method chaining
         */
        public Builder addGlobalCacheKeys(List<String> globalCacheKeys) {
            params.globalCacheKeys.addAll(globalCacheKeys);
            return this;
        }

        /**
         * Add a global cache key for global cache request
         *
         * @param globalCacheKey The global cache key to add
         * @return This builder for method chaining
         */
        public Builder addGlobalCacheKey(String globalCacheKey) {
            params.globalCacheKeys.add(globalCacheKey);
            return this;
        }

        /**
         * Build the FeatureFlowParameters object with the provided parameters.
         *
         * @return The fully constructed FeatureFlowParameters instance
         */
        public FeatureFlowParameters build() {
            return params;
        }
    }
}
