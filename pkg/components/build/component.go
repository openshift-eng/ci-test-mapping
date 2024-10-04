package build

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var BuildComponent = Component{
	Component: &config.Component{
		Name:                 "Build",
		Operators:            []string{},
		DefaultJiraComponent: "Build",
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-builds",
			},
			{
				SIG: "sig-devex",
			},
			{
				IncludeAny: []string{
					":BuildAPI ",
				},
				Priority: 1,
			},
			{Suite: "BuildAPI"},
			{Suite: "build related"},
			{Suite: "template related scenarios:"},
			{Suite: "Testing the isolation during build scenarios"},
			{Suite: "buildlogic.feature"},
			{Suite: "admin build related features"},
		},
		TestRenames: map[string]string{
			"[sig-devex][Feature:Templates] templateinstance readiness test should report failed soon after an annotated objects has failed [apigroup:template.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]":                                                            "[sig-devex][Feature:Templates] templateinstance readiness test should report failed soon after an annotated objects has failed [apigroup:template.openshift.io][apigroup:build.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
			"[sig-devex][Feature:Templates] templateinstance readiness test should report ready soon after all annotated objects are ready [apigroup:template.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]":                                                             "[sig-devex][Feature:Templates] templateinstance readiness test should report ready soon after all annotated objects are ready [apigroup:template.openshift.io][apigroup:build.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
			"[sig-devex][Feature:ImageEcosystem][perl][Slow] hot deploy for openshift perl image hot deploy test should work [apigroup:image.openshift.io][apigroup:operator.openshift.io][apigroup:config.openshift.io][apigroup:build.openshift.io]":                                                                                                          "[sig-devex][Feature:ImageEcosystem][perl][Slow] hot deploy for openshift perl image hot deploy test should work [apigroup:image.openshift.io][apigroup:operator.openshift.io][apigroup:config.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io]",
			"[sig-devex][Feature:Templates] templateinstance readiness test should report ready soon after all annotated objects are ready [apigroup:template.openshift.io][apigroup:build.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]":                                                                                         "[sig-devex][Feature:Templates] templateinstance readiness test should report ready soon after all annotated objects are ready [apigroup:template.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
			"[sig-builds][Feature:Builds][volumes] build volumes should mount given secrets and configmaps into the build pod for docker strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io] [Suite:openshift/conformance/parallel]":                                                                                                     "[sig-builds][Feature:Builds][volumes] build volumes should mount given secrets and configmaps into the build pod for docker strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-builds][Feature:Builds][volumes] build volumes should mount given secrets and configmaps into the build pod for source strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io] [Suite:openshift/conformance/parallel]":                                                                                                     "[sig-builds][Feature:Builds][volumes] build volumes should mount given secrets and configmaps into the build pod for source strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret into the build pod for docker strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io] [Suite:openshift/conformance/parallel]":                          "[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret into the build pod for docker strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret into the build pod for source strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io] [Suite:openshift/conformance/parallel]":                          "[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret into the build pod for source strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret without resource refresh into the build pod for docker strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io] [Suite:openshift/conformance/parallel]": "[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret without resource refresh into the build pod for docker strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret without resource refresh into the build pod for source strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io] [Suite:openshift/conformance/parallel]": "[sig-builds][Feature:Builds][volumes] csi build volumes within Tech Preview enabled cluster [apigroup:config.openshift.io] should mount given csi shared resource secret without resource refresh into the build pod for source strategy builds [apigroup:image.openshift.io][apigroup:build.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/parallel]",
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
