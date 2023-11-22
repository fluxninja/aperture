import React, { useEffect, useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { dark } from "react-syntax-highlighter/dist/esm/styles/prism";

const CodeSnippet = ({ lang, snippetName }) => {
  const [code, setCode] = useState("");

  useEffect(() => {
    const fetchSnippet = async () => {
      try {
        const response = await import("./code-snippets.json"); // Update the path to your JSON file
        setCode(response.default[lang][snippetName]);
      } catch (error) {
        console.error("Error fetching snippet:", error);
      }
    };

    fetchSnippet();
  }, [lang, snippetName]);

  return (
    <SyntaxHighlighter language={lang} style={dark}>
      {code}
    </SyntaxHighlighter>
  );
};

export default CodeSnippet;
