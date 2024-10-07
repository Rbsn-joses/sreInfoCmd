package pbi

import (
	"context"
	"strings"

	"github.com/Rbsn-joses/create-pbi/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/sirupsen/logrus"
)

var (
	wiqlQueryGetPBI = "SELECT * FROM WorkItems WHERE [System.WorkItemType] ='Product Backlog Item' AND [System.AreaPath] = 'SRE'"
	Logger          *logrus.Logger
)

func CheckIfProductBacklogItemExist(productBacklogItemName string, pbiNames []types.PBI_DEVOPS) string {
	var pbiURL string
	for _, pbi := range pbiNames {
		if strings.Contains(pbi.Title, productBacklogItemName) {
			pbiURL = pbi.URL
			break
		}
	}

	return pbiURL
}
func GetProductBacklogItem(witClient workitemtracking.Client, ctx context.Context) []types.PBI_DEVOPS {

	//project := "teste-api"
	//team := "teste-api Team"
	// Create query arguments
	queryArgs := workitemtracking.QueryByWiqlArgs{
		Wiql: &workitemtracking.Wiql{
			Query: &wiqlQueryGetPBI,
		},
	}

	// Get first page of the list of team projects for your organization
	responseValue, err := witClient.QueryByWiql(ctx, queryArgs)
	if err != nil {
		Logger.Error(err)
	}
	pbiInfo := getProductBacklogItemInfo(witClient, ctx, responseValue)
	return pbiInfo

}
func getProductBacklogItemInfo(witClient workitemtracking.Client, ctx context.Context, responseValue *workitemtracking.WorkItemQueryResult) []types.PBI_DEVOPS {
	var PBIs []types.PBI_DEVOPS
	var PBI types.PBI_DEVOPS

	for _, info := range *responseValue.WorkItems {
		responsePBIInfo, err := witClient.GetWorkItem(ctx, workitemtracking.GetWorkItemArgs{Id: info.Id})
		if err != nil {
			Logger.Error(err)
		}
		PBI.Title = (*responsePBIInfo.Fields)["System.Title"].(string)
		PBI.URL = *responsePBIInfo.Url
		PBIs = append(PBIs, PBI)
	}
	return PBIs
}
