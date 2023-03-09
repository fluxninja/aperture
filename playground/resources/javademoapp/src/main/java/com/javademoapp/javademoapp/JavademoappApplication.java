package com.javademoapp.javademoapp;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.propagation.TextMapPropagator;
import io.opentelemetry.context.propagation.TextMapSetter;
import org.apache.catalina.LifecycleException;
import org.apache.catalina.connector.Response;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.net.HttpURLConnection;
import java.net.URLConnection;
import java.time.Duration;

@SpringBootApplication
@RestController
public class JavademoappApplication {

	private static final String LIBRARY_NAME = "javademoapp";
	public static void main(String[] args) {
		SpringApplication.run(JavademoappApplication.class, args);
	}

	@GetMapping("/")
	public String demoApp() throws LifecycleException, ServletException {

		String hostname = getHostnamefromEnv();
		int port = getPortfromEnv();
		int concurrency = getConcurrencyfromEnv();
		Duration latency = getLatencyfromEnv();
		double rejectRatio = getRejectRatiofromEnv();

		SimpleService simpleService = new SimpleService(hostname, port, concurrency, latency, rejectRatio);
		HttpServletRequest request = null;
		HttpServletResponse response = new Response();
		simpleService.run(request, response);
		return "Demo App running successfully";
	}

	public String getHostnamefromEnv() {
		if (System.getenv("HOSTNAME") == null)
			return "localhost";
		else
			return System.getenv("HOSTNAME");
	}

	public int getPortfromEnv() {
		if (System.getenv("SIMPLE_SERVICE_PORT") == null)
			return 8088;
		else return Integer.parseInt(System.getenv("SIMPLE_SERVICE_PORT"));
	}

	public int getConcurrencyfromEnv() {
		String concurrency = System.getenv("SIMPLE_SERVICE_CONCURRENCY");
		if (concurrency == null) {
			return 10;
		}
		return Integer.parseInt(concurrency);
	}

	public Duration getLatencyfromEnv(){
		String latency = System.getenv("SIMPLE_SERVICE_LATENCY");
		if (latency == null) {
			return Duration.ofMillis(50);
		}
		return Duration.parse(latency);
	}

	public double getRejectRatiofromEnv(){
		String rejectRatio = System.getenv("SIMPLE_SERVICE_REJECT_RATIO");
		if (rejectRatio == null) {
			return 0.05;
		}
		return Double.parseDouble(rejectRatio);
	}

	@Bean
	public Tracer tracer() {
		return OpenTelemetry.noop().getTracer(LIBRARY_NAME);
	}

	@Bean
	public TextMapPropagator textMapPropagator() {
		return OpenTelemetry.noop().getPropagators().getTextMapPropagator();
	}

	@Bean
	public TextMapSetter<HttpURLConnection> textMapSetter() {
		return URLConnection::setRequestProperty;
	}
}
