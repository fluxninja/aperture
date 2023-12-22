import { ApertureClient } from "@fluxninja/aperture-js";

async function initializeApertureClient() {
    const address = process.env.APERTURE_AGENT_ADDRESS || 'localhost:8089';
    const apiKey = process.env.APERTURE_API_KEY || '';

    const apertureClient = new ApertureClient({
        address: address,
        apiKey: apiKey,
    });

    return apertureClient;
}

async function sendRequest(apertureClient: ApertureClient) {
    const flow = await apertureClient.startFlow("concurrency-limiter-feature", {
        labels: {
            user_id: "some_user_id",
        },
        grpcCallOptions: {
            deadline: Date.now() + 300,
        },
    });

    if (flow.shouldRun()) {
        console.log("Request accepted. Processing...");
    } else {
        console.log("Request rejected due to concurrency limit. Try again later.");
    }

    flow.end();
}

async function handleConcurrencyLimit(apertureClient: ApertureClient) {
    const requestsPerSecond = 10;
    const durationInSeconds = 200;

    for (let i = 0; i < durationInSeconds; i++) {
        const requests = Array.from({ length: requestsPerSecond }, () => sendRequest(apertureClient));
        // sending requests in parallel to simulate concurrency
        await Promise.all(requests);

        // Wait 1 second before sending the next batch of requests
        await new Promise(resolve => setTimeout(resolve, 1000));
    }
}

async function main() {
    const apertureClient = await initializeApertureClient();
    await handleConcurrencyLimit(apertureClient);
}

main().catch(e => console.error(e));
