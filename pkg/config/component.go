package config

import (
	"regexp"
	"strings"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Component is the default configuration struct that you can include in your
// own component implementation. It includes a matcher help that will identify
// if a test belongs to a sig, operator, as well as simple substring matching.
// Components do not need to use this framework, it's an optional add-on.
type Component struct {
	Name                 string
	DefaultJiraComponent string
	Matchers             []ComponentMatcher
	Operators            []string
	Namespaces           []string

	// When a test is renamed, you can still look at results across releases by mapping new names
	// to the oldest version of the test.
	TestRenames map[string]string
}

// ComponentMatcher is used to match against a TestInfo struct. Note the fields SIG,
// Suite, IncludeAll, and ExcludeAll are ANDed together. That is, all that have values must
// match.  For include  and exclude, the individual items in the array are ANDed. That
// is, if you  specify multiple substrings, all must match. Use separate component
// matchers for an OR operation.
//
// The second set  of fields are metadata used to assign ownership.
type ComponentMatcher struct {
	SIG        string
	Suite      string
	IncludeAll []string
	IncludeAny []string
	ExcludeAll []string
	ExcludeAny []string

	JiraComponent string
	Capabilities  []string
	Priority      int
}

func (c *Component) FindMatch(test *v1.TestInfo) *ComponentMatcher {
	if ok, capabilities := c.IsOperatorTest(test); ok {
		return &ComponentMatcher{
			JiraComponent: c.DefaultJiraComponent,
			Capabilities:  capabilities,
		}
	}

	if namespace, ok := c.IsNamespaceTest(test.Name); ok {
		if c.IsInNamespace(namespace) {
			return &ComponentMatcher{
				JiraComponent: c.DefaultJiraComponent,
			}
		}
		return nil
	}

	// Check if any of the Matchers match the given test
	for _, m := range c.Matchers {
		sigMatch := true
		suiteMatch := true
		namespaceMatch := true
		incSubstrMatch := true
		incAnySubstrMatch := true

		if m.SIG != "" {
			sigMatch = util.IsSigTest(test.Name, m.SIG)
		}

		if m.Suite != "" {
			suiteMatch = m.IsSuiteTest(test)
		}

		if len(m.IncludeAll) > 0 {
			incSubstrMatch = m.IsSubstringAllTest(m.IncludeAll, test)
		}
		if len(m.IncludeAny) > 0 {
			incAnySubstrMatch = m.IsSubstringAnyTest(m.IncludeAny, test)
		}

		if len(m.ExcludeAll) > 0 {
			// If all the exclusions are present, we force a non-match
			if m.IsSubstringAllTest(m.ExcludeAll, test) {
				continue
			}
		}
		if len(m.ExcludeAny) > 0 {
			// If any of the exclusions are present, we force a non-match
			if m.IsSubstringAnyTest(m.ExcludeAny, test) {
				continue
			}
		}

		// AND the match results together
		if sigMatch && suiteMatch && namespaceMatch && incSubstrMatch && incAnySubstrMatch {
			return &m
		}
	}

	return nil
}

func (c *Component) ListNamespaces() []string {
	return sets.NewString(c.Namespaces...).List()
}

func (c *Component) IsInNamespace(testNamespace string) bool {
	for _, namespace := range c.Namespaces {
		if testNamespace == namespace {
			return true
		}
	}
	return false
}

func (c *Component) IsNamespaceTest(testName string) (string, bool) {
	testNamespace := ExtractNamespaceFromTestName(testName)
	return testNamespace, len(testNamespace) > 0
}

func (cm *ComponentMatcher) IsSuiteTest(test *v1.TestInfo) bool {
	return test.Suite == cm.Suite
}

func (cm *ComponentMatcher) IsSubstringAllTest(allOf []string, test *v1.TestInfo) bool {
	for _, str := range allOf {
		if !strings.Contains(test.Name, str) {
			return false
		}
	}
	return true
}

func (cm *ComponentMatcher) IsSubstringAnyTest(anyOf []string, test *v1.TestInfo) bool {
	for _, str := range anyOf {
		if strings.Contains(test.Name, str) {
			return true
		}
	}
	return false
}

func (c *Component) IsOperatorTest(test *v1.TestInfo) (bool, []string) {
	for _, operator := range c.Operators {
		// OpenShift tests related to operators (install, upgrade, etc)
		if isOperatorTest, capabilities := util.IdentifyOperatorTest(operator, test.Name); isOperatorTest {
			return true, capabilities
		}
	}

	return false, nil
}

var namespaceShort = regexp.MustCompile(`ns/(?P<Namespace>[-\w]+)`)
var namespaceFull = regexp.MustCompile(`namespace/(?P<Namespace>[-\w]+)`)

func ExtractNamespaceFromTestName(in string) string {
	if namespaceShort.MatchString(in) {
		return namespaceShort.FindStringSubmatch(in)[1]
	}
	if namespaceFull.MatchString(in) {
		return namespaceFull.FindStringSubmatch(in)[1]
	}
	return ""
}
