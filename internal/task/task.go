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
	ID       string    `yaml:"id" json:"id"`
	Title    string    `yaml:"title" json:"title"`
	Status   string    `yaml:"status" json:"status"`
	Priority string    `yaml:"priority" json:"priority"`
	Assignee string    `yaml:"assignee,omitempty" json:"assignee"`
	Tags     FlowSlice `yaml:"tags,omitempty" json:"tags"`
	Links    FlowSlice `yaml:"links,omitempty" json:"links"`
	Created  string    `yaml:"created" json:"created"`
	Updated  string    `yaml:"updated" json:"updated"`
	Body     string    `yaml:"-" json:"body"`
}
