package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.FlowControlServiceGrpc;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.google.gson.Gson;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.time.Duration;
import java.util.concurrent.TimeUnit;

class FlowControlClient {
    private final FlowControlServiceGrpc.FlowControlServiceBlockingStub grpcFlowControlClient;
    private final boolean forceHttp;
    private final String agentHost;
    private final String agentPort;

    FlowControlClient(
            boolean forceHttp,
            String agentHost,
            String agentPort,
            FlowControlServiceGrpc.FlowControlServiceBlockingStub grpcFlowControlClient) {
        this.forceHttp = forceHttp;
        this.agentHost = agentHost;
        this.agentPort = agentPort;
        this.grpcFlowControlClient = grpcFlowControlClient;
    }

    // TODO: add http timeout
    CheckResponse check(CheckRequest req, Duration timeout) {
        if(this.forceHttp) {
            String requestGson = new Gson().toJson(req);
            System.out.println(requestGson);

            String checkUrl = agentHost + ":" + agentPort + "/Check";
            try {
                URL url = new URL(checkUrl);
                HttpURLConnection con = (HttpURLConnection) url.openConnection();

                con.setRequestMethod("POST");
                con.setRequestProperty("Content-Type", "application/json");

                con.setDoOutput(true);
                OutputStream os = con.getOutputStream();
                os.write(requestGson.getBytes());
                os.flush();
                os.close();

                int responseCode = con.getResponseCode();
                System.out.println("RESPONSE: " + responseCode);

                BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream()));
                String inputLine;
                StringBuffer response = new StringBuffer();
                while ((inputLine = in.readLine()) != null) {
                    response.append(inputLine);
                }
                in.close();

                String jsonResp = response.toString();
                System.out.println("Response Body : " + jsonResp);

                Gson responseGson = new Gson();
                CheckResponse responseObject = responseGson.fromJson(jsonResp, CheckResponse.class);
                System.out.println("Parsed Response Object : " + responseObject);
                return responseObject;
            } catch (Exception e) {
                // TODO: proper exceptions
                System.out.println("oops");
                return null;
            }
        } else {
            return this.grpcFlowControlClient
                    .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
                    .check(req);
        }
    }
}
