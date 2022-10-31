export const FlowStatus = Object.freeze({
    Ok: Symbol(0),
    Error: Symbol(1)
})

export class Flow {
    constructor(checkResponse) {
        this.checkResponse = checkResponse;
    }

    Accepted() {
        return true;
    }

    End(flowStatus) {
        return true;
    }

    CheckResponse() {
        return this.checkResponse;
    }
}
