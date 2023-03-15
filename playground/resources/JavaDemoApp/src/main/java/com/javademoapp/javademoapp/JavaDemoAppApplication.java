package com.javademoapp.javademoapp;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.time.Duration;

@SpringBootApplication
public class JavaDemoAppApplication {

	public static final String DEFAULT_APP_PORT = "8080";
	public static final String DEFAULT_AGENT_HOST = "localhost";
	public static final String DEFAULT_AGENT_PORT = "8089";

	public static final String DEFAULT_CONCURRENCY = "10";
	public static final String DEFAULT_LATENCY = "50";
	public static final String DEFAULT_REJECT_RATIO = "0.05";

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

		String hostname = System.getenv("HOSTNAME");
		if (hostname == null) {
			hostname = DEFAULT_AGENT_HOST;
		}
		System.setProperty("HOSTNAME", hostname);

		String port = System.getenv("PORT");
		if (port == null) {
			port = DEFAULT_APP_PORT;
		}
		System.setProperty("PORT", port);

		String concurrency = System.getenv("CONCURRENCY");
		if (concurrency == null) {
			concurrency = DEFAULT_CONCURRENCY;
		}
		System.setProperty("CONCURRENCY", concurrency);

		String latency = System.getenv("LATENCY");
		if (latency == null) {
			latency = DEFAULT_LATENCY;
		}
		System.setProperty("LATENCY", latency);

		String rejectRatio = System.getenv("REJECT_RATIO");
		if (rejectRatio == null) {
			rejectRatio = DEFAULT_REJECT_RATIO;
		}
		System.setProperty("REJECT_RATIO", rejectRatio);

		SpringApplication.run(JavaDemoAppApplication.class, args);
	}

}
