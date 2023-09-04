export type DecisionType = "DECISION_TYPE_ACCEPTED" | string;

export type Response = {
  decisionType: DecisionType;
};

export type Error = { code: number };
