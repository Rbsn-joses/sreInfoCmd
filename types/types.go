package types

import "github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"

type PBI_EXCEL struct {
	Titulo      string
	Name        string
	Data        string
	Tags        string
	Description string
}
type PBI_DEVOPS struct {
	Title string `json:"title"`
	URL   string
}

type ReposInfo struct {
	Name            string
	VariableGroups  []taskagent.VariableGroup
	RepoDescription string
}
type VariableGroup struct {
	Name                     string
	VariableGroupDescription string
	Variables                map[string]any
}
