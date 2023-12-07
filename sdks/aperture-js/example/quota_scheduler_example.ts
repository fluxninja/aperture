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

// START: Priority
const userTiers = {
    "platinum": 8,
    "gold": 4,
    "silver": 2,
    "free": 1,
};
// END: Priority

const intervalTime = 1000;

async function sendRequestForTier(apertureClient: ApertureClient, tier: string, priority: number) {
    // START: QSStartFlow
    const flow = await apertureClient.startFlow("my-feature", {
        labels: {
            user_id: "some_user_id",
            priority: priority.toString(),
            workload: `${tier} user`,
        },
        grpcCallOptions: {
            deadline: Date.now() + 120000, // ms
        },
    });

    console.log(`Request sent for ${tier} tier with priority ${priority}.`);
    flow.end();
    // END: QSStartFlow
}

function scheduleRequests(apertureClient: ApertureClient) {
    Object.entries(userTiers).forEach(([tier, priority]) => {
        setInterval(() => {
            sendRequestForTier(apertureClient, tier, priority);
        }, intervalTime);
    });
}

async function main() {
    const apertureClient = await initializeApertureClient();
    scheduleRequests(apertureClient);
}

main().catch(e => console.error(e));
