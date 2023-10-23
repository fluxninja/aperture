// BlueprintsComponents.js
import React from "react";
import CodeBlock from "@theme/CodeBlock";
import TabItem from "@theme/TabItem";

export const TabContent = ({ tabValue, valuesFile, policyName }) => {
  let commands = '';

  switch (tabValue) {
    case "aperturectl (Aperture Cloud)":
      commands = `aperturectl cloud blueprints apply --values-file=${valuesFile}.yaml`;
      break;
    case "aperturectl (self-hosted controller)":
      commands = `
aperturectl blueprints generate --values-file=${valuesFile}.yaml --output-dir=policy-gen
aperturectl apply policy --file=policy-gen/policies/${policyName}.yaml --kube
      `;
      break;
    case "kubectl (self-hosted controller)":
      commands = `
aperturectl blueprints generate --values-file=${valuesFile}.yaml --output-dir=policy-gen
kubectl apply -f policy-gen/policies/${policyName}-cr.yaml -n aperture-controller
      `;
      break;
    default:
      return null;
  }

  return (
    <CodeBlock language="bash">
      {commands}
    </CodeBlock>
  );
};

export const BashTab = ({ tabValue, children }) => (
  <TabItem value={tabValue} label={tabValue}>
    {children}
  </TabItem>
);
