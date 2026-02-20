package task

import (
	"gopkg.in/yaml.v3"
)

// FlowSlice marshals as a YAML flow sequence: [a, b, c]
type FlowSlice []string

func (f FlowSlice) MarshalYAML() (any, error) {
	node := &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FlowStyle,
	}
	for _, s := range f {
		node.Content = append(node.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: s,
		})
	}
	return node, nil
}

type Task struct {
	ID       string    `yaml:"id"`
	Title    string    `yaml:"title"`
	Status   string    `yaml:"status"`
	Priority string    `yaml:"priority"`
	Assignee string    `yaml:"assignee,omitempty"`
	Tags     FlowSlice `yaml:"tags,omitempty"`
	Links    FlowSlice `yaml:"links,omitempty"`
	Created  string    `yaml:"created"`
	Updated  string    `yaml:"updated"`
	Body     string    `yaml:"-"`
}
