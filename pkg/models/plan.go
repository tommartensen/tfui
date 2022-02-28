package models

type change struct {
	Actions      []string               `json:"actions"`
	Before       map[string]interface{} `json:"before"`
	After        map[string]interface{} `json:"after"`
	AfterUnknown map[string]interface{} `json:"after_unknown"`
}

type resourceChange struct {
	Address      string `json:"address"`
	Mode         string `json:"mode"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Change       change `json:"change"`
}

type TfPlanMeta struct {
	Region    string `json:"region"`
	Project   string `json:"project"`
	Workspace string `json:"workspace"`
	Date      string `json:"date"`
	CommitID  string `json:"commit_id"`
}

type TfPlan struct {
	FormatVersion   string           `json:"format_version"`
	TfVersion       string           `json:"terraform_version"`
	ResourceChanges []resourceChange `json:"resource_changes"`
	Meta            TfPlanMeta       `json:"meta"`
}

type PlannedChanges struct {
	Create int `json:"create"`
	Delete int `json:"delete"`
	NoOp   int `json:"no-op"`
	Update int `json:"update"`
}

type TfPlanSummary struct {
	Meta           TfPlanMeta     `json:"meta"`
	PlannedChanges PlannedChanges `json:"plannedChanges"`
	State          string         `json:"state"`
}
