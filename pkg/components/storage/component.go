package storage

import (
	"regexp"

	log "github.com/sirupsen/logrus"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

// A regular expression + its replacement string to apply to test names
type testNameReplacement struct {
	regexp      *regexp.Regexp
	replacement string
}

type Component struct {
	*config.Component

	testNameReplacements []testNameReplacement
}

var StorageComponent = Component{
	Component: &config.Component{
		Name:                 "Storage",
		Operators:            []string{"storage", "csi-snapshot-controller"},
		DefaultJiraComponent: "Storage",
		Namespaces: []string{
			"openshift-cluster-csi-drivers",
			"openshift-cluster-storage-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				SIG:      "sig-storage",
				Priority: 1,
			},
			{
				IncludeAny: []string{
					"bz-Storage",
					"bz-storage",
				},
			},
			{
				IncludeAny: []string{
					":Storage ",
					"upgrade should succeed: storage",
					"upgrade should succeed: csi-snapshot-controller",
				},
				Priority: 1,
			},
			{
				IncludeAny: []string{
					"[k8s.io] GKE local SSD [Feature:GKELocalSSD] should write and read from node local SSD [Feature:GKELocalSSD] [sig-arch] [Suite:openshift/conformance/parallel] [Suite:k8s]",
					"[k8s.io] GKE node pools [Feature:GKENodePool] should create a cluster with multiple node pools [Feature:GKENodePool] [sig-arch] [Suite:openshift/conformance/parallel] [Suite:k8s]",
				},
			},
			{Suite: "All in one volume"},
			{Suite: "Azure disk and Azure file specific scenarios"},
			{Suite: "CSI Resizing related feature"},
			{Suite: "CSI snapshot operator related scenarios"},
			{Suite: "CSI testing related feature"},
			{Suite: "STORAGE"},
			{Suite: "Storage upgrade tests"},
			{Suite: "cluster storage operator related scenarios"},
			{Suite: "storage security check"},
			{Suite: `"storage (storageclass, pv, pvc) related"`},
			{Suite: "storage security check"},
			{Suite: "storageClass related feature"},
			{Suite: "NFS Persistent Volume"},
			{Suite: "Persistent Volume reclaim policy tests"},
			{Suite: "ResourceQuata for storage"},
			{Suite: "PVC resizing Test"},
			{Suite: "storage (storageclass, pv, pvc) related"},
			{Suite: "Scenarios specific for block volume support"},
			{Suite: "CSI snapshot webhook related scenarios"},
			{Suite: "testing for parameter fsType"},
			{Suite: "Persistent Volume Claim binding policies"},
			{Suite: "Storage of Hostpath plugin testing"},
			{Suite: "Testing for pv and pvc pre-bind feature"},
			{Suite: "Target pvc to a specific pv"},
		},
		TestRenames: map[string]string{
			"[Storage][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-csi-drivers":         "[bz-Storage][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-csi-drivers",
			"[Storage][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-storage-operator":    "[bz-Storage][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-storage-operator",
			"[Storage][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-csi-drivers":      "[bz-Storage][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-csi-drivers",
			"[Storage][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-storage-operator": "[bz-Storage][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-storage-operator",
		},
	},

	testNameReplacements: []testNameReplacement{
		{
			// Remove [MinimumKubeletVersion:1.27] introduced in Kubernetes 1.29 read-write-once-pod tests to match 1.28 test names.
			regexp:      regexp.MustCompile(`\[MinimumKubeletVersion:1\.27\]`),
			replacement: "",
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
	// Apply any test name replacements
	name := test.Name
	for _, r := range c.testNameReplacements {
		name = r.regexp.ReplaceAllString(name, r.replacement)
	}

	// Look up the stable name for our test in our renamed tests map.
	if stableName, ok := c.TestRenames[name]; ok {
		name = stableName
	}
	if name != test.Name {
		log.Tracef("Mapped storage test %q to %q", test.Name, name)
	}
	return name
}

func (c *Component) JiraComponents() (components []string) {
	components = []string{c.DefaultJiraComponent}
	for _, m := range c.Matchers {
		components = append(components, m.JiraComponent)
	}

	return components
}
