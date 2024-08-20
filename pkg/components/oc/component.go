package oc

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OcComponent = Component{
	Component: &config.Component{
		Name:                 "oc",
		Operators:            []string{},
		DefaultJiraComponent: "oc",
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-cli",
			},
			{
				IncludeAny: []string{"[sig-arch] External binary usage"},
			},
			{Suite: "Workloads"},
			{Suite: `"Check status via oc status, wait etc"`},
			{Suite: "Update emptyDir volumes via command: oc set volumes"},
			{Suite: "basic verification for upgrade oc client testing"},
			{Suite: "build 'apps' with CLI"},
			{Suite: "creating 'apps' with CLI"},
			{Suite: "oc annotate related features"},
			{Suite: "oc credential related scenarios"},
			{Suite: "oc debug related scenarios"},
			{Suite: "oc explain resources for storage"},
			{Suite: "oc get related command"},
			{Suite: "oc inspect related scenarios"},
			{Suite: "oc logs related features"},
			{Suite: "oc new-app related scenarios"},
			{Suite: "oc observe related tests"},
			{Suite: "oc patch/apply related scenarios"},
			{Suite: "oc related features"},
			{Suite: "oc run related scenarios"},
			{Suite: "oc set image related tests"},
			{Suite: "oc set triggers tests"},
			{Suite: "oc tag related scenarios"},
			{Suite: "oc_delete.feature"},
			{Suite: "oc_expose.feature"},
			{Suite: "oc_login.feature"},
			{Suite: "oc_secrets.feature"},
			{Suite: "oc_set_probe.feature"},
			{Suite: "oc_set_resources.feature"},
			{Suite: "oc_volume.feature"},
			{Suite: "projects related features via cli"},
			{Suite: "route related features via cli"},
			{Suite: "set deployment-hook/build-hook with CLI"},
			{Suite: "oc image mirror related scenarios"},
			{Suite: "oc_rsync.feature"},
			{Suite: "oc_portforward.feature"},
			{Suite: "oc_process.feature"},
			{Suite: "Add, update remove volume to rc/dc and --overwrite option"},
			{Suite: "rsh.feature"},
			{Suite: "oc import-image related feature"},
			{Suite: "admin deployment related features"},
			{Suite: "Check status via oc status, wait etc"},
			{Suite: "oc_set_env.feature"},
			{Suite: "oc idle"},
			{Suite: "oc extract related scenarios"},
			{Suite: "Return description with cli"},
			{Suite: "oc proxy related scenarios"},
			{Suite: "oc plugin related tests"},
			{Suite: "oc/kubernetes version related features"},
		},
		TestRenames: map[string]string{
			"[sig-cli] oc adm build-chain [apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:project.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/parallel]":                        "[sig-cli] oc adm build-chain [apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:project.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-cli] oc builds complex build start-build [apigroup:build.openshift.io][apigroup:apps.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]":                                            "[sig-cli] oc builds complex build start-build [apigroup:build.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
			"[sig-cli] oc builds complex build webhooks CRUD [apigroup:build.openshift.io][apigroup:apps.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]":                                          "[sig-cli] oc builds complex build webhooks CRUD [apigroup:build.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
			"[sig-cli] oc debug deployment from a build [apigroup:image.openshift.io] [Suite:openshift/conformance/parallel]":                                                                                                  "[sig-cli] oc debug deployment configs from a build [apigroup:image.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-cli] oc status returns expected help messages [apigroup:project.openshift.io][apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:route.openshift.io] [Suite:openshift/conformance/parallel]": "[sig-cli] oc status returns expected help messages [apigroup:project.openshift.io][apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:apps.openshift.io][apigroup:route.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-cli] templates process [apigroup:template.openshift.io][Skipped:Disconnected] [Suite:openshift/conformance/parallel]":                                                                                        "[sig-cli] templates process [apigroup:apps.openshift.io][apigroup:template.openshift.io][Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
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
