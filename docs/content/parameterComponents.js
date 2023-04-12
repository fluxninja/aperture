// ParameterComponents.js

import React from "react";

export const ParameterHeading = ({ children }) => (
  <span style={{ fontWeight: "bold" }}>{children}</span>
);

export const WrappedDescription = ({ children }) => (
  <span style={{ wordWrap: "normal" }}>{children}</span>
);

export const RefType = ({ type, reference }) => <a href={reference}>{type}</a>;

export const ParameterDescription = ({
  name,
  type,
  reference,
  value,
  description,
}) => (
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
      <td>
        <code>{value}</code>
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
  </table>
);
