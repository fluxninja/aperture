import React from "react";

const VersionCode = (props) => {
  const childrenAsString = React.Children.toArray(props.children).join("");
  return (
    <pre>
      <code>{childrenAsString.replace("${version}", props.version)}</code>
    </pre>
  );
};

export default VersionCode;
