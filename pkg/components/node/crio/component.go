package nodecrio

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var CRIOComponent = Component{
	Component: &config.Component{
		Name:                 "Node / CRI-O",
		Operators:            []string{},
		DefaultJiraComponent: "Node / CRI-O",
		TestRenames: map[string]string{
			// SigstoreImageVerification tests
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerification][Serial] Should fail clusterimagepolicy signature validation root of trust does not match the identity in the signature":                  "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerification][Serial] Should fail clusterimagepolicy signature validation root of trust does not match the identity in the signature [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerification][Serial] Should fail clusterimagepolicy signature validation when scope in allowedRegistries list does not skip signature verification":   "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerification][Serial] Should fail clusterimagepolicy signature validation when scope in allowedRegistries list does not skip signature verification [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerification][Serial] Should fail imagepolicy signature validation in different namespaces root of trust does not match the identity in the signature": "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerification][Serial] Should fail imagepolicy signature validation in different namespaces root of trust does not match the identity in the signature [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerification][Serial] Should pass clusterimagepolicy signature validation with signed image":                                                           "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerification][Serial] Should pass clusterimagepolicy signature validation with signed image [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerification][Serial] Should pass imagepolicy signature validation with signed image in namespaces":                                                    "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerification][Serial] Should pass imagepolicy signature validation with signed image in namespaces [Suite:openshift/conformance/serial]",
			// SigstoreImageVerificationPKI tests
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] clusterimagepolicy signature validation tests fail with PKI email does not match":                                       "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] clusterimagepolicy signature validation tests fail with PKI email does not match [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] clusterimagepolicy signature validation tests fail with PKI root of trust does not match the identity in the signature": "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] clusterimagepolicy signature validation tests fail with PKI root of trust does not match the identity in the signature [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] clusterimagepolicy signature validation tests pass with valid PKI":                                                      "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] clusterimagepolicy signature validation tests pass with valid PKI [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] imagepolicy signature validation tests fail with PKI root of trust does not match the identity in the signature":        "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] imagepolicy signature validation tests fail with PKI root of trust does not match the identity in the signature [Suite:openshift/conformance/serial]",
			"[sig-imagepolicy][Suite:openshift/disruptive-longrunning][Disruptive][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] imagepolicy signature validation tests pass with valid PKI":                                                             "[sig-imagepolicy][OCPFeatureGate:SigstoreImageVerificationPKI][Serial][Skipped:Disconnected] imagepolicy signature validation tests pass with valid PKI [Suite:openshift/conformance/serial]",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"[sig-arch] [Conformance] FIPS TestFIPS",
				},
			},
			{Suite: "Container_Engine_Tools"},
			{
				SIG: "sig-imagepolicy",
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
