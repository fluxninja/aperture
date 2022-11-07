import {
    FEATURE_STATUS_LABEL,
    CHECK_RESPONSE_LABEL,
    FLOW_END_TIMESTAMP_LABEL,
} from "./consts.js";

export const FlowStatus = Object.freeze({
    Ok: "Ok",
    Error: "Error"
});

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
        if (this.checkResponse.decisionType === 'DECISION_TYPE_ACCEPTED') {
            return true;
        }
        return false;
    }

    End(flowStatus) {
        if (this.ended) {
            return new Error("flow already ended");
        }
        this.ended = true;

        this.span.setAttribute(FEATURE_STATUS_LABEL, flowStatus);
        this.span.setAttribute(CHECK_RESPONSE_LABEL, JSON.stringify(this.checkResponse));
        this.span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now());

        let attr = this.span.attributes;
        this.span.end();
    }

    CheckResponse() {
        return this.checkResponse;
    }
}
