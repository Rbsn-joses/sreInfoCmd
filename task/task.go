package task

import (
	"context"
	"fmt"

	"github.com/Rbsn-joses/create-pbi/excel"
	"github.com/Rbsn-joses/create-pbi/pbi"
	"github.com/Rbsn-joses/create-pbi/types"
	"github.com/sirupsen/logrus"

	"github.com/microsoft/azure-devops-go-api/azuredevops"

	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
)

var (
	OrganizationUrl     string
	PersonalAccessToken string
	Project             string
	ExcelFile           string
	TypeWorkItem        = "Task"
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

	connection := azuredevops.NewPatConnection(organizationUrl, PAT)
	ctx := context.Background()
	// Create a work item tracking client
	witClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		Logger.Errorf("Error creating work item tracking client: %v", err)
		return
	}
	err = checkTasksExist(witClient, ctx, UsernameDevops)
	if err != nil {
		Logger.Error(err)
	}
}

func checkTasksExist(witClient workitemtracking.Client, ctx context.Context, userName string) error {
	var wiqlQueryGetTasks = fmt.Sprintf("SELECT * FROM WorkItems WHERE [System.WorkItemType] ='Task' AND [System.State] = 'To Do' AND [System.AssignedTo] = '%s'", userName)

	queryArgs := workitemtracking.QueryByWiqlArgs{
		Wiql: &workitemtracking.Wiql{
			Query: &wiqlQueryGetTasks,
		},
	}

	// Get first page of the list of team projects for your organization
	responseValue, err := witClient.QueryByWiql(ctx, queryArgs)
	if err != nil {
		return err
	}
	taskLists := getNewTaskInfoByUser(witClient, ctx, responseValue)
	exceltable, f, err := excel.GetExcel(Logger, ExcelFile)
	if err != nil {
		return err
	}

	newTasks := findUniqueValues(taskLists, exceltable)
	pbisDevops := pbi.GetProductBacklogItem(witClient, ctx)
	createdUrlWorkItemSuccess := createWorkitemTask(witClient, ctx, newTasks, pbisDevops)
	excel.FinalizeAndSaveExcel(f, createdUrlWorkItemSuccess)
	return nil
}

func findUniqueValues(taskLists map[string]string, exceltable []types.PBI_EXCEL) []types.PBI_EXCEL {

	// Find the unique structs based on the "Titulo" field
	uniqueStructs := []types.PBI_EXCEL{}
	for _, s := range exceltable {
		if taskLists[s.Titulo] == "" {
			uniqueStructs = append(uniqueStructs, s)
		}
	}
	Logger.Debug("task únicas no arquivo excel", uniqueStructs)
	return uniqueStructs
}

func getNewTaskInfoByUser(witClient workitemtracking.Client, ctx context.Context, responseValue *workitemtracking.WorkItemQueryResult) map[string]string {
	// Get first page of the list of team projects for your organization
	var listTaskTitle = make(map[string]string)
	for _, workItem := range *responseValue.WorkItems {
		responseWorkItemByID, err := witClient.GetWorkItem(ctx, workitemtracking.GetWorkItemArgs{Id: workItem.Id})
		if err != nil {
			Logger.Fatal(err)
		}
		title := (*responseWorkItemByID.Fields)["System.Title"].(string)
		listTaskTitle[title] = title

	}
	Logger.Debug("task puxadas no azure devops", listTaskTitle)

	return listTaskTitle

}

func createWorkitemTask(witClient workitemtracking.Client, ctx context.Context, newTasks []types.PBI_EXCEL, pbiDevops []types.PBI_DEVOPS) []*workitemtracking.WorkItem {
	pathTitle := "/fields/System.Title"
	pathRelation := "/relations/-"
	pathUserAssined := "/fields/System.AssignedTo"
	pathWorkType := "/fields/System.WorkItemType"
	pathTags := "/fields/System.Tags"
	var createdWorkItemList []*workitemtracking.WorkItem
	for _, task := range newTasks {
		ProductBacklogItemParentLink := pbi.CheckIfProductBacklogItemExist(task.Name, pbiDevops)
		if ProductBacklogItemParentLink != "" {
			workitem := []webapi.JsonPatchOperation{
				{
					Op:    &webapi.OperationValues.Add,
					Path:  &pathTitle,
					Value: task.Titulo,
				},
				{
					Op:   &webapi.OperationValues.Add,
					Path: &pathRelation,
					Value: map[string]interface{}{
						"rel": "System.LinkTypes.Hierarchy-Reverse",
						"url": ProductBacklogItemParentLink,
					},
				},
				{
					Op:    &webapi.OperationValues.Add,
					Path:  &pathUserAssined,
					Value: UserName,
				},
				{
					Op:    &webapi.OperationValues.Add,
					Path:  &pathWorkType,
					Value: "Task",
				},
				{
					Op:    &webapi.OperationValues.Add,
					Path:  &pathTags,
					Value: task.Tags,
				},
			}
			// Create the work item
			createdWorkItem, err := witClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{Document: &workitem, Project: &Project, Type: &TypeWorkItem})
			if err != nil {
				Logger.Errorf("Error creating work item: %v", err)
				return nil
			}

			Logger.Infof(fmt.Sprintf("Nova Task adicionada com sucesso url: %s, titulo: %s", *createdWorkItem.Url, (*createdWorkItem.Fields)["System.Title"].(string)))
			createdWorkItemList = append(createdWorkItemList, createdWorkItem)

		} else {
			Logger.Warningf(fmt.Sprintf("task não criada devido a inexistência do productBacklogItem definito no excel titulo: %s data: %s PBI: %s", task.Titulo, task.Data, task.Name))
		}
	}
	return createdWorkItemList

}
