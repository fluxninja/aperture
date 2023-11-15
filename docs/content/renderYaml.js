import React, { useState, useEffect } from 'react';

const RenderYamlWithoutUri = ({ yamlPath }) => {
  const [yamlContent, setYamlContent] = useState('');

  useEffect(() => {
    fetch(yamlPath)
      .then(response => response.text())
      .then(data => {
        const processedData = data
          .split('\n')
          .filter(line => !line.trim().startsWith('uri'))
          .join('\n');
        setYamlContent(processedData);
      })
      .catch(error => console.error('Error loading YAML file:', error));
  }, [yamlPath]);

  return <pre>{yamlContent}</pre>;
};

export default RenderYamlWithoutUri;
