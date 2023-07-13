package storage

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var StorageComponent = Component{
	Component: &config.Component{
		Name:                 "Storage",
		Operators:            []string{"storage", "csi-snapshot-controller"},
		DefaultJiraComponent: "Storage",
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-storage",
			},
			{
				IncludeAny: []string{
					"bz-Storage",
					"bz-storage",
				},
			},
			{
				Namespaces: []string{
					"openshift-cluster-csi-drivers",
					"openshift-cluster-storage-operator",
				},
			},
			{
				IncludeAny: []string{
					"[k8s.io] GKE local SSD [Feature:GKELocalSSD] should write and read from node local SSD [Feature:GKELocalSSD] [sig-arch] [Suite:openshift/conformance/parallel] [Suite:k8s]",
					"[k8s.io] GKE node pools [Feature:GKENodePool] should create a cluster with multiple node pools [Feature:GKENodePool] [sig-arch] [Suite:openshift/conformance/parallel] [Suite:k8s]",
				},
			},
		},
	},
}

func (c *Component) IdentifyTest(test *v1.TestInfo) (*v1.TestOwnership, error) {
	if matcher := c.FindMatch(test); matcher != nil {
		jira := matcher.JiraComponent
		if jira == "" {
			jira = c.DefaultJiraComponent
		}
		return &v1.TestOwnership{
			Name:          test.Name,
			Component:     c.Name,
			JIRAComponent: jira,
			Priority:      matcher.Priority,
			Capabilities:  append(matcher.Capabilities, identifyCapabilities(test)...),
		}, nil
	}

	return nil, nil
}

func (c *Component) StableID(test *v1.TestInfo) string {
	// Look up the stable name for our test in our renamed tests map.
	if stableName, ok := c.TestRenames[test.Name]; ok {
		return stableName
	}
	return test.Name
}

func (c *Component) JiraComponents() (components []string) {
	components = []string{c.DefaultJiraComponent}
	for _, m := range c.Matchers {
		components = append(components, m.JiraComponent)
	}

	return components
}
