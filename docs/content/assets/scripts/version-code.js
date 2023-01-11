import React, { useRef } from "react";

const VersionCode = (props) => {
  const codeRef = useRef(null);

  const handleCopyClick = () => {
    const text = codeRef.current.textContent;
    navigator.clipboard.writeText(text);
  };

  const childrenAsString = React.Children.toArray(props.children).join("");
  return (
    <div>
      <pre className="highlight" ref={codeRef}>
        <code>{childrenAsString.replace("${version}", props.version)}</code>
      </pre>
      <button className="copy-button" onClick={handleCopyClick}>
        Copy
      </button>
    </div>
  );
};

export default VersionCode;
