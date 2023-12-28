import { ApertureClient } from "@fluxninja/aperture-js";
import { FlowStatus, LookupStatus } from "@fluxninja/aperture-js";
import { Request, Response } from "express";

// Create aperture client
export const apertureClient = new ApertureClient({
  address: "ORGANIZATION.app.fluxninja.com:443",
  apiKey: "API_KEY",
});


async function UIRLLabelMatcher(apertureClient: ApertureClient){
    // START: UIRLLabelMatcher
    const flow = await apertureClient.startFlow("rate-limiting-feature", {
        labels: {
            user_id: "user1",
            customer_tier: "gold",
            product_tier: "trial",
          },
        grpcCallOptions: {
            deadline: Date.now() + 300, // ms
        },
    });
    // END: UIRLLabelMatcher

}

async function QSUI(apertureClient: ApertureClient, tier: string, priority: number) {
    // START: QSUI
    const flow = await apertureClient.startFlow("quota-scheduling-feature", {
        labels: {
            user_id: "some_user_id",
            priority: "100",
            workload: "gold user",
        },
        grpcCallOptions: {
            deadline: Date.now() + 120000, // ms
        },
    });
    console.log(`Request sent for ${tier} tier with priority ${priority}.`);
    flow.end();
    // END: QSUI
}

async function UIRLTokens(apertureClient: ApertureClient){
    // START: UIRLTokens
    const flow = await apertureClient.startFlow("rate-limiting-feature", {
        labels: {
          user_id: "user1",
          tier: "premium",
          tokens: "50",

        },
        grpcCallOptions: {
          deadline: Date.now() + 300, // 300ms deadline
        },
      });
    // END: UIRLTokens
}

const userTiers: { [key: string]: number } = {
    "platinum": 8,
    "gold": 4,
    "silver": 2,
    "free": 1,
};

async function QSPriority(apertureClient: ApertureClient, tier: string, priority: number) {
    // START: QSPriority
    const userPriority = userTiers[tier] || 1;

    const flow = await apertureClient.startFlow("quota-scheduling-feature", {
        labels: {
            user_id: "some_user_id",
            priority: userPriority.toString(),
            workload: `${tier} user`,
        },
        grpcCallOptions: {
            deadline: Date.now() + 120000, // ms
        },
    });
    console.log(`Request sent for ${tier} tier with priority ${userPriority}.`);
    flow.end();
    // END: QSPriority
}

async function UIQSTokens(apertureClient: ApertureClient, tier: string, priority: number, userType: string) {
    // START: UIQSTokens
    let userTokens;
    switch (userType) {
    case "premium":
        userTokens = 100;
        break;
    case "gold":
        userTokens = 50;
        break;
    default:
        userTokens = 0;
    }
    const flow = await apertureClient.startFlow("quota-scheduling-feature", {
        labels: {
            user_id: "some_user_id",
            product_tier: "trial",
            priority: "100",
            tokens: userTokens.toString(),
        },
        grpcCallOptions: {
            deadline: Date.now() + 120000, // ms
        },
    });
    // END: UIQSTokens
}

async function UIWorkload(apertureClient: ApertureClient, tier: string, priority: number) {
    // START: UIQSWorkload
    let userWorkload = "subscriber";
    const flow = await apertureClient.startFlow("quota-scheduling-feature", {
        labels: {
            user_id: "some_user_id",
            product_tier: "trial",
            priority: priority.toString(),
            workload: userWorkload,
        },
        grpcCallOptions: {
            deadline: Date.now() + 120000, // ms
        },
    });
    // END: UIQSWorkload
}

async function UIConcurrencyTokens(apertureClient: ApertureClient, tier: string, priority: number) {
    // START: UIConcurrencyTokens
    const flow = await apertureClient.startFlow("concurrency-limiting-feature", {
        labels: {
          user_id: "user1",
          tier: "premium",
          tokens: "50",

        },
        grpcCallOptions: {
          deadline: Date.now() + 300, // 300ms deadline
        },
      });
    // END: UIConcurrencyTokens
}
