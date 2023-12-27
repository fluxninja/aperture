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
  const flow = await apertureClient.startFlow("concurrency-limiting-feature", {
    labels: {
      user_id: "some_user_id",
    },
    grpcCallOptions: {
      deadline: Date.now() + 300,
    },
  });

  if (flow.shouldRun()) {
    //console.log("Request accepted. Processing..." + flow.checkResponse());
  } else {
    console.log("Request rejected due to concurrency limit. Try again later.");
  }
  // do a json stringify for checkResponse() to work
  console.log(JSON.stringify(flow.checkResponse()));

  flow.end();
}

async function handleConcurrencyLimit(apertureClient: ApertureClient) {
  const requestsPerSecond = 10;
  const durationInSeconds = 200;

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
