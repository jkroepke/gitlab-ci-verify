package format_conversion

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseYamlNode(t *testing.T) {
	_, err := ParseYamlNode([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseYamlNodeComponent(t *testing.T) {
	content := []byte(`spec:
  inputs:
    image:
      type: string
---
deploy:
  script:
    - echo $UNQUOTED_VARIABLE
`)
	doc, err := ParseYamlNode(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		t.Fatalf("expected a combined mapping document, got %#v", doc)
	}
	var foundSpec, foundDeploy bool
	for idx := 0; idx < len(doc.Content[0].Content); idx += 2 {
		key := doc.Content[0].Content[idx]
		switch key.Value {
		case "spec":
			foundSpec = true
		case "deploy":
			foundDeploy = true
			if key.Line != 6 {
				t.Fatalf("expected deploy source line 6, got %d", key.Line)
			}
		}
	}
	if !foundSpec || !foundDeploy {
		t.Fatalf("expected both spec and deploy in merged document")
	}
}
