package v1

type JiraUser struct {
	Self        string `json:"self"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	TimeZone    string `json:"timeZone"`
}

type JiraComponent struct {
	Self                string   `json:"self"`
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Lead                JiraUser `json:"lead"`
	AssigneeType        string   `json:"assigneeType"`
	Assignee            JiraUser `json:"assignee"`
	RealAssigneeType    string   `json:"realAssigneeType"`
	RealAssignee        JiraUser `json:"realAssignee"`
	IsAssigneeTypeValid bool     `json:"isAssigneeTypeValid"`
	Project             string   `json:"project"`
	ProjectID           int      `json:"projectId"`
	Archived            bool     `json:"archived"`
	Deleted             bool     `json:"deleted"`
}
