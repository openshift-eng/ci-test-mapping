package v1

type Config struct {
	// IncludeSuites is a specific list of suites to include from the tests table.
	IncludeSuites []string `yaml:"includeSuites"`

	// IncludeSuitePatterns is a list of SQL LIKE patterns for suite names to include (e.g. "test-suite%").
	IncludeSuitePatterns []string `yaml:"includeSuitePatterns"`

	// ExcludeSuites is a specific list of suite names to exclude (exact match).
	ExcludeSuites []string `yaml:"excludeSuites"`

	// ExcludeSuitePatterns is a list of SQL LIKE patterns for suite names to exclude.
	ExcludeSuitePatterns []string `yaml:"excludeSuitePatterns"`

	// ExcludeTests is a specific list of tests to exclude from the tests table.
	ExcludeTests []string `yaml:"excludeTests"`

	// IncludeJobs is a specific list of CI jobs to include from the tests table.
	IncludeJobs []string `yaml:"includeJobs"`

	// ExcludeJobs is a specific list of CI jobs to exclude from the tests table.
	ExcludeJobs []string `yaml:"excludeJobs"`
}
