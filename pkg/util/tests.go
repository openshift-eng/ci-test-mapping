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

// ExtractTestField gets the value of a field in a test name. Fields are formatted either was [Field: Value]
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

// StableID produces a stable test ID based on a TestInfo struct and a stableName.
func StableID(testInfo v1.TestInfo, stableName string) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(stableName)))
	if testInfo.Suite != "" {
		stableName = fmt.Sprintf("%s:%s", testInfo.Suite, hash)
	}

	return stableName
}
