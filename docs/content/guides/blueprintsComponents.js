// BlueprintsComponents.js

import React from "react";
import CodeBlock from "@theme/CodeBlock";
import TabItem from "@theme/TabItem";

export const TabContent = ({ tabValue, valuesFile, policyName }) => {

  const commandMap = {
    "aperturectl (Aperture Cloud)": `aperturectl cloud blueprints apply --values-file=${valuesFile}.yaml`,

    "aperturectl (self-hosted controller)": `
aperturectl blueprints generate --values-file=${valuesFile}.yaml --output-dir=policy-gen
aperturectl apply policy --file=policy-gen/policies/${policyName}.yaml --kube
    `,

    "kubectl (self-hosted controller)": `
aperturectl blueprints generate --values-file=${valuesFile}.yaml --output-dir=policy-gen
kubectl apply -f policy-gen/policies/${policyName}-cr.yaml -n aperture-controller
    `
  };

  const commands = commandMap[tabValue];
  if (!commands) return null;

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
