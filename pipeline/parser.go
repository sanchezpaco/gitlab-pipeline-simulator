package pipeline

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func ExpandYAML(data []byte) *yaml.Node {
	root := expandYAML(data)
	return root
}

func expandYAML(data []byte) *yaml.Node {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		panic(fmt.Errorf("failed to unmarshal YAML: %w", err))
	}

	anchors := make(map[string]*yaml.Node)
	collectAnchors(&root, anchors)
	resolveAliases(&root, anchors)
	processMergeKeys(&root, anchors)

	if len(root.Content) > 0 {
		docRoot := root.Content[0]
		resolveReferences(docRoot, docRoot)
	}

	return &root
}

func extractJobsAndStages(root *yaml.Node) ([]Job, []string, error) {
	var jobs []Job
	var stageOrder []string

	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, nil, fmt.Errorf("invalid YAML structure")
	}

	content := root.Content[0]
	if content.Kind != yaml.MappingNode {
		return nil, nil, fmt.Errorf("root node is not a mapping")
	}

	// Extract stage order
	for i := 0; i < len(content.Content); i += 2 {
		key := content.Content[i]
		value := content.Content[i+1]

		if key.Value == "stages" && value.Kind == yaml.SequenceNode {
			for _, stageNode := range value.Content {
				if stageNode.Kind == yaml.ScalarNode {
					stageOrder = append(stageOrder, stageNode.Value)
				}
			}
		}
	}

	// Extract jobs
	for i := 0; i < len(content.Content); i += 2 {
		key := content.Content[i]
		value := content.Content[i+1]

		if key.Kind == yaml.ScalarNode && key.Value != "stages" && value.Kind == yaml.MappingNode {
			job := Job{Name: key.Value}
			for j := 0; j < len(value.Content); j += 2 {
				attrKey := value.Content[j]
				attrValue := value.Content[j+1]

				switch attrKey.Value {
				case "stage":
					if attrValue.Kind == yaml.ScalarNode {
						job.Stage = attrValue.Value
					}
				case "rules":
					if attrValue.Kind == yaml.SequenceNode {
						for _, ruleNode := range attrValue.Content {
							if ruleNode.Kind == yaml.MappingNode {
								rule := Rule{}
								for k := 0; k < len(ruleNode.Content); k += 2 {
									ruleKey := ruleNode.Content[k]
									ruleValue := ruleNode.Content[k+1]

									switch ruleKey.Value {
									case "if":
										if ruleValue.Kind == yaml.ScalarNode {
											rule.If = ruleValue.Value
										}
									case "when":
										if ruleValue.Kind == yaml.ScalarNode {
											rule.When = ruleValue.Value
										}
									}
								}
								job.Rules = append(job.Rules, rule)
							}
						}
					}
				case "script":
					if attrValue.Kind == yaml.SequenceNode {
						for _, scriptNode := range attrValue.Content {
							if scriptNode.Kind == yaml.ScalarNode {
								job.Scripts = append(job.Scripts, scriptNode.Value)
							}
						}
					}
				}
			}
			jobs = append(jobs, job)
		}
	}

	return jobs, stageOrder, nil
}

func collectAnchors(node *yaml.Node, anchors map[string]*yaml.Node) {
	if node == nil {
		return
	}

	if node.Anchor != "" {
		anchors[node.Anchor] = node
	}

	for _, child := range node.Content {
		collectAnchors(child, anchors)
	}
}

func resolveAliases(node *yaml.Node, anchors map[string]*yaml.Node) {
	if node == nil {
		return
	}

	if node.Kind == yaml.AliasNode {
		if anchored, ok := anchors[node.Alias.Anchor]; ok {
			*node = *anchored
		}
	}

	for _, child := range node.Content {
		resolveAliases(child, anchors)
	}
}

func processMergeKeys(node *yaml.Node, anchors map[string]*yaml.Node) {
	if node.Kind != yaml.MappingNode {
		for _, child := range node.Content {
			processMergeKeys(child, anchors)
		}
		return
	}

	var newContent []*yaml.Node
	var merges []*yaml.Node

	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		if key.Value == "<<" {
			merges = append(merges, getMergeSources(value, anchors)...)
		} else {
			newContent = append(newContent, key, value)
		}
	}

	for _, merge := range merges {
		if merge.Kind != yaml.MappingNode {
			continue
		}
		newContent = append(merge.Content, newContent...)
	}

	node.Content = newContent

	for _, child := range node.Content {
		processMergeKeys(child, anchors)
	}
}

func getMergeSources(node *yaml.Node, anchors map[string]*yaml.Node) []*yaml.Node {
	var sources []*yaml.Node

	if node.Kind == yaml.SequenceNode {
		for _, item := range node.Content {
			resolveAliases(item, anchors)
			sources = append(sources, item)
		}
	} else {
		resolveAliases(node, anchors)
		sources = append(sources, node)
	}

	return sources
}

func resolveReferences(node *yaml.Node, root *yaml.Node) {
	if node == nil {
		return
	}

	if node.Tag == "!reference" {
		if path := parseReferencePath(node); len(path) > 0 {
			if referenced, err := findReferencedNode(root, path); err == nil {
				*node = *referenced
			}
		}
	}

	for _, child := range node.Content {
		resolveReferences(child, root)
	}
}

func parseReferencePath(node *yaml.Node) []string {
	var path []string

	if node.Kind == yaml.SequenceNode {
		for _, item := range node.Content {
			if item.Kind == yaml.ScalarNode {
				path = append(path, item.Value)
			}
		}
	}

	return path
}

func findReferencedNode(root *yaml.Node, path []string) (*yaml.Node, error) {
	current := root

	for _, component := range path {
		if current.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("nodo no es un mapping en el componente: %s", component)
		}

		found := false
		for i := 0; i < len(current.Content); i += 2 {
			key := current.Content[i]
			if key.Value == component {
				current = current.Content[i+1]
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("componente no encontrado: %s", component)
		}
	}

	return current, nil
}
