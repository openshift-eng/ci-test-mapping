package jira

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

func GetJiraComponents() (map[string]int64, error) {
	start := time.Now()
	log.Infof("loading jira ocpbugs component information...")
	body, err := jiraRequest("https://issues.redhat.com/rest/api/2/project/12332330/components")
	if err != nil {
		return nil, err
	}

	var components []v1.JiraComponent
	err = json.Unmarshal(body, &components)
	if err != nil {
		return nil, errors.WithMessage(err, "response from jira: "+string(body))
	}

	ids := make(map[string]int64)
	for _, c := range components {
		jiraID, err := strconv.ParseInt(c.ID, 10, 64)
		if err != nil {
			log.WithError(err).Warnf("error parsing jira ID '%s'", c.ID)
			continue
		}

		ids[c.Name] = jiraID
	}

	log.Infof("jira ocpbugs components loaded in %+v", time.Since(start))
	return ids, nil
}

func jiraRequest(apiURL string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
