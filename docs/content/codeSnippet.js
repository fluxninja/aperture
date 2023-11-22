import React, { useEffect, useState } from "react";
import CodeBlock from "@theme/CodeBlock";

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

  return <CodeBlock language={lang}>{code}</CodeBlock>;
};

export default CodeSnippet;
