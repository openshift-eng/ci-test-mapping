package tracinguiplugin

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
)

func identifyCapabilities(test *v1.TestInfo) []string {
	capabilities := util.DefaultCapabilities(test)

	// Extract [Capability:XXX] tags from test name
	capabilities = append(capabilities, util.ExtractTestField(test.Name, "Capability")...)

	return capabilities
}
