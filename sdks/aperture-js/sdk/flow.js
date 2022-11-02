import { fcs } from "./utils.js";
import {
    FEATURE_STATUS_LABEL,
    CHECK_RESPONSE_LABEL,
    FLOW_END_TIMESTAMP_LABEL,
} from "./consts.js";

export const FlowStatus = Object.freeze({
    Ok: Symbol(0),
    Error: Symbol(1)
})

export class Flow {
    constructor(span, checkResponse = null) {
        this.span = span;
        this.ended = false;
        this.checkResponse = checkResponse;
    }

    Accepted() {
        if (this.checkResponse === undefined) {
            return true;
        }
        if (this.checkResponse.DecisionType === fcs.CheckResponse_DECISION_TYPE_ACCEPTED) {
            return true;
        }
        return false;
    }

    End(flowStatus) {
        if (this.ended) {
            return new Error("flow already ended");
        }
        this.ended = true;

        this.span.setAttributes({
            FEATURE_STATUS_LABEL: flowStatus,
            CHECK_RESPONSE_LABEL: this.checkResponse,
            FLOW_END_TIMESTAMP_LABEL: Date.now(),
        });
        this.span.end();
    }

    CheckResponse() {
        return this.checkResponse;
    }
}
