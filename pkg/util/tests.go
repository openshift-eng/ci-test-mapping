package util

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

var fieldRegexp = regexp.MustCompile(`\[([^\]]*):([^\]]*)\]`)

// ExtractTestField gets the value of a field in a test name, e.g. [sig-storage][Driver: gce] would return
// "gce" when extracting "Driver"
func ExtractTestField(testName, field string) (results []string) {
	matches := fieldRegexp.FindAllStringSubmatch(testName, -1)
	for _, match := range matches {
		if len(match) == 3 && match[1] == field {
			results = append(results, strings.TrimSpace(match[2]))
		}
	}

	return results
}

// StableID produces a stable test ID based on a TestInfo struct and a stableName.
func StableID(testInfo *v1.TestInfo, stableName string) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(stableName)))
	if testInfo.Suite != "" {
		// We differentiate tests on suite and name, because I imagined a scenario where two entirely different
		// test suites create a generically named test like "TestClusterCreate" and we don't want to conflate them.
		// However, openshift-tests and  openshift-tests-upgrade are definitely the same tests, and it makes sense to
		// give them both the same ID. We already get the "upgrade" view from the variants.
		suite := testInfo.Suite
		if suite == "openshift-tests-upgrade" {
			suite = "openshift-tests"
		}
		stableName = fmt.Sprintf("%s:%s", suite, hash)
	}

	return stableName
}
