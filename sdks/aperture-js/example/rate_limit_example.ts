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

async function handleRequestRateLimit(apertureClient: ApertureClient) {
    while (true) {
        // START: RLStartFlow
        const flow = await apertureClient.startFlow("rate-limiting-feature", {
            labels: {
                user_id: "some_user_id",
            },
            grpcCallOptions: {
                deadline: Date.now() + 300, // ms
            },
        });
        // END: RLStartFlow

        await new Promise(resolve => setTimeout(resolve, 1000));

        // START: RLFlowShouldRun
        if (flow.shouldRun()) {
            console.log("Request accepted. Processing...");
        } else {
            console.log("Request rate-limited. Try again later.");
        }

        flow.end();
        // END: RLFlowShouldRun
    }
}

async function main() {
    const apertureClient = await initializeApertureClient();
    await handleRequestRateLimit(apertureClient);
}

main().catch(e => console.error(e));
