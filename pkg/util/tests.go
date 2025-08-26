package util

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

var fieldRegexp = regexp.MustCompile(`(\[([^\]]*):([^\]]*)\]|(\w+)\/("([^"]*)"|\S+))`)

// ExtractTestField gets the value of a field in a test name. Fields are formatted either as [Field: Value]
// or Field/Value.  Field is case-insensitive.
func ExtractTestField(testName, field string) (results []string) {
	matches := fieldRegexp.FindAllStringSubmatch(testName, -1)
	for _, match := range matches {
		count := len(match)
		for i, matchField := range match {
			if !strings.EqualFold(matchField, field) || count < i+2 {
				continue
			}

			value := strings.TrimSpace(match[i+1])
			unquoted, err := strconv.Unquote(value)
			if err == nil {
				value = unquoted
			}
			results = append(results, value)
		}
	}

	return results
}

// removeTestField looks for matching fields and removes
// the complete field & value pair from the testName
func removeTestField(testName, field string) string {
	updatedTestName := testName
	matches := fieldRegexp.FindAllStringSubmatch(testName, -1)
	for _, match := range matches {
		if len(match) > 0 {
			// we want the first match which should be the complete field value
			matched := match[0]
			if !strings.Contains(matched, field) {
				continue
			}
			// does the match have a leading " ", if so remove it too
			index := strings.Index(updatedTestName, matched)
			if index > 0 && updatedTestName[index-1] == ' ' {
				matched = fmt.Sprintf(" %s", matched)
			}

			updatedTestName = strings.TrimSpace(strings.Replace(updatedTestName, matched, "", 1))
			// we found our field so quit processing
			break
		}
	}

	return updatedTestName
}

// StableID produces a stable test ID based on a TestInfo struct and a stableName.
func StableID(testInfo *v1.TestInfo, stableName string) string {
	// Monitor attribute will be added automatically to monitor tests
	// want to avoid mass mapping so we remove from the stableID
	stableName = removeTestField(stableName, "[Monitor:")
	hash := fmt.Sprintf("%x", md5.Sum([]byte(stableName)))
	if testInfo.Suite != "" {
		stableName = fmt.Sprintf("%s:%s", testInfo.Suite, hash)
	}

	return stableName
}
