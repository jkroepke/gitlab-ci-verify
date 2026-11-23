package format_conversion

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// ParseYamlNode parses all documents in a GitLab CI YAML stream.
// CI/CD components use a spec header separated from the job definitions by ---.
// Keep the original nodes and line numbers for reporting findings in source files.
func ParseYamlNode(content []byte) (*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	decoder.UniqueKeys(false)

	var result yaml.Node
	for {
		var document yaml.Node
		err := decoder.Decode(&document)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(document.Content) == 0 {
			continue
		}
		if document.Content[0].Kind != yaml.MappingNode {
			return nil, fmt.Errorf("expected mapping in GitLab CI YAML document")
		}
		if len(result.Content) == 0 {
			result = document
			continue
		}
		result.Content[0].Content = append(result.Content[0].Content, document.Content[0].Content...)
	}

	return &result, nil
}
