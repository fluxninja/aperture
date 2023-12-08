import { ApertureClient, LookupStatus } from "@fluxninja/aperture-js";
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

const intervalTime = 1000;

async function monitorCacheAndUpdate(apertureClient: ApertureClient) {
    while (true) {
        // START: CStartFlow
        const flow = await apertureClient.startFlow("my-feature", {
            labels: {
                user_id: "some_user_id",
            },
            grpcCallOptions: {
                deadline: Date.now() + 120000, // ms
            },
            resultCacheKey: "cache",
        });
        // END: CStartFlow

        await new Promise(resolve => setTimeout(resolve, intervalTime));

        // START: CFlowShouldRun
        if (flow.shouldRun()) {
            if (flow.resultCache().getLookupStatus() === LookupStatus.Hit) {
                console.log(flow.resultCache().getValue()?.toString());
            } else {
                console.log("Cache miss, setting cache value");

                const resString = "foo";
                const buffer = Buffer.from(resString);
                const setResp = await flow.setResultCache({
                    value: buffer,
                    ttl: {
                        seconds: 10,
                        nanos: 0,
                    },
                });
                if (setResp.getError()) {
                    console.log(`Error setting cache value: ${setResp.getError()}`);
                }
                console.log(resString);
            }
        }

        flow.end();
        // END: CFlowShouldRun
    }
}

async function main() {
    const apertureClient = await initializeApertureClient();
    monitorCacheAndUpdate(apertureClient);
}

main().catch(e => console.error(e));
