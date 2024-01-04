import { ApertureClient } from "@fluxninja/aperture-js";
import inquirer from 'inquirer';

async function initializeApertureClient() {
    const answers = await inquirer.prompt([
        {
            type: 'input',
            name: 'address',
            message: 'Enter your organization\'s address:',
        },
        {
            type: 'input',
            name: 'apiKey',
            message: 'Enter the API key:',
        },
    ]);

    const apertureClient = new ApertureClient({
        address: answers.address,
        apiKey: answers.apiKey,
    });

    return apertureClient;
}

async function sendRequest(apertureClient: ApertureClient) {
  // START: CLStartFlow
  const flow = await apertureClient.startFlow("concurrency-limiting-feature", {
    labels: {
      user_id: "some_user_id",
    },
    grpcCallOptions: {
      deadline: Date.now() + 1000,
    },
  });
  // END: CLStartFlow

  // START: CLFlowShouldRun
  if (flow.shouldRun()) {
    console.log("Request accepted. Processing..." + flow.checkResponse());
  } else {
    console.log("Request rejected due to concurrency limit. Try again later.");
  }

  await flow.end();
  // END: CLFlowShouldRun
}

async function handleConcurrencyLimit(apertureClient: ApertureClient) {
  const requestsPerSecond = 5;
  const durationInSeconds = 1000;
  for (let i = 0; i < durationInSeconds; i++) {
    const requests = Array.from({ length: requestsPerSecond }, () =>
      sendRequest(apertureClient),
    );
    // sending requests in parallel to simulate concurrency
    await Promise.all(requests);

    // Wait 1 second before sending the next batch of requests
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
}

async function main() {
  const apertureClient = await initializeApertureClient();
  await handleConcurrencyLimit(apertureClient);
}

main().catch((e) => console.error(e));
