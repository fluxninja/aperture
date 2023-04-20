// ParameterComponents.js

import CodeBlock from "@theme/CodeBlock";
import yaml from "js-yaml";
import React from "react";

export const ParameterHeading = ({ children }) => (
  <span style={{ fontWeight: "bold" }}>{children}</span>
);

export const WrappedDescription = ({ children }) => (
  <span style={{ wordWrap: "normal" }}>{children}</span>
);

export const RefType = ({ type, reference }) => <a href={reference}>{type}</a>;

const renderSimpleValue = (value) => <code>{value}</code>;

const renderComplexValue = (yamlValue) => (
  <details>
    <summary>Expand</summary>
    <CodeBlock language="yaml">{yamlValue}</CodeBlock>
  </details>
);

const renderValue = (jsonValue) => {
  if (
    typeof jsonValue === "string" ||
    typeof jsonValue === "number" ||
    typeof jsonValue === "boolean" ||
    jsonValue === null
  ) {
    return renderSimpleValue(jsonValue);
  } else {
    const yamlValue = yaml.dump(jsonValue);
    return renderComplexValue(yamlValue);
  }
};

export const ParameterDescription = ({
  name,
  type,
  reference,
  value,
  description,
}) => {
  const jsonValue = JSON.parse(value);

  return (
    <table className="blueprints-params">
      <tr>
        <td>
          <ParameterHeading>Parameter</ParameterHeading>
        </td>
        <td>
          <code>{name}</code>
        </td>
      </tr>
      <tr>
        <td className="blueprints-description">
          <ParameterHeading>Description</ParameterHeading>
        </td>
        <td className="blueprints-description">
          <WrappedDescription>{description}</WrappedDescription>
        </td>
      </tr>
      <tr>
        <td>
          <ParameterHeading>Type</ParameterHeading>
        </td>
        <td>
          <em>
            {reference === "" ? (
              type
            ) : (
              <RefType type={type} reference={reference} />
            )}
          </em>
        </td>
      </tr>
      <tr>
        <td className="blueprints-default-heading">
          <ParameterHeading>Default Value</ParameterHeading>
        </td>
        <td>{renderValue(jsonValue)}</td>
      </tr>
    </table>
  );
};
