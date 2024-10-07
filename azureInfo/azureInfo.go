package azureinfo

import (
	"context"
	"fmt"
	"log"

	"github.com/Rbsn-joses/create-pbi/excel"
	"github.com/Rbsn-joses/create-pbi/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"

	"github.com/sirupsen/logrus"
)

var (
	OrganizationUrl     string
	PersonalAccessToken string
	Project             string
	ExcelFile           string
	UserName            string
	Logger              *logrus.Logger
)

func StartInjection(UsernameDevops, PAT, project, organizationUrl, excelFile string, logger *logrus.Logger) {
	Logger = logger
	OrganizationUrl = organizationUrl
	Project = project
	PersonalAccessToken = PAT
	ExcelFile = excelFile
	UserName = UsernameDevops
	var ReposInfo types.ReposInfo
	var ReposInfoList []types.ReposInfo
	connection := azuredevops.NewPatConnection(organizationUrl, PAT)
	ctx := context.Background()
	taskagentclient, err := taskagent.NewClient(ctx, connection)
	if err != nil {
		fmt.Println(err)
	}
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		fmt.Println(err)
	}
	responseValue, err := coreClient.GetProjects(ctx, core.GetProjectsArgs{})
	if err != nil {
		log.Fatal(err)
	}
	//Logger.Debug("len ", len((*responseValue).Value))

	for _, project := range (*responseValue).Value {
		ReposInfo.Name = *project.Name
		ReposInfo.RepoDescription = *project.Description
		varGroup, err := taskagentclient.GetVariableGroups(ctx, taskagent.GetVariableGroupsArgs{Project: project.Name})
		if err != nil {
			fmt.Println(err)
		}
		ReposInfo.VariableGroups = *varGroup

		ReposInfoList = append(ReposInfoList, ReposInfo)

	}
	//Logger.Debug(ReposInfoList)

	excel.CreateExcel(ReposInfoList, logger)

}
func checkStrucIsEmpty(variableGroup taskagent.VariableGroup) bool {
	return variableGroup == taskagent.VariableGroup{}
}
