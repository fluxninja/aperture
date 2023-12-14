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
            user_agent: "user_agent1",
            customer_tier: "gold",
            product_tier: "trial",
          },
        grpcCallOptions: {
            deadline: Date.now() + 300, // ms
        },
    });
    // END: UIRLLabelMatcher

}

async function UIQSBlueprint(req: Request, res: Response) {
    const userTiers = {
        "platinum": 8,
        "gold": 4,
        "silver": 2,
        "free": 1,
    };
    // START: UIQSBlueprint
    const flow = await apertureClient.startFlow("quota-scheduling-feature", {
        labels: {
            user_id: "some_user_id",
            priority: "8",
            workload: `platinum`,
        },
        grpcCallOptions: {
            deadline: Date.now() + 120000, // ms
        },
    });
    // END: UIQSBlueprint
}

async function UIRLTokens(apertureClient: ApertureClient){
    // START: UIRLTokens
    const flow = await apertureClient.startFlow("rate-limiting-feature", {
        labels: {
          limit_key: "user1",
          tier: "premium",
          tokens: "50",

        },
        grpcCallOptions: {
          deadline: Date.now() + 300, // 300ms deadline
        },
      });
    // END: UIRLTokens
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
            priority: priority.toString(),
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
