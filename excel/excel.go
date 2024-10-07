package excel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Rbsn-joses/create-pbi/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

var Logger *logrus.Logger
var cabeçalhosEsperados = []string{"titulo", "tags", "data", "pbi", "description"}
var sheetName = strings.ToLower("relatorio-pbi")
var f *excelize.File
var err error

func GetExcel(logger *logrus.Logger, ExcelFile string) ([]types.PBI_EXCEL, *excelize.File, error) {
	Logger = logger
	var pbis []types.PBI_EXCEL
	var pbi types.PBI_EXCEL

	// Abre o arquivo Excel
	f, err = excelize.OpenFile(ExcelFile)
	if err != nil {
		Logger.Error(err)
		return nil, nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		Logger.Error(err)
		return nil, nil, err
	}
	headers := rows[0]
	excelFormatIsCorrect := checkHeadersExists(headers, cabeçalhosEsperados)
	if !excelFormatIsCorrect {
		errmessage := fmt.Sprintf("formato do arquivo %s no sheetname %s está com formatação incorreta de headers\n", ExcelFile, sheetName)
		errmessage = errmessage + fmt.Sprintf("verificar ordem dos headers e nomeclatura para seguir com o processo, headers esperados: %v\n", strings.ToLower(strings.Join(cabeçalhosEsperados, ",")))
		errmessage = errmessage + fmt.Sprintf("headers no excel: %v", strings.ToLower(strings.Join(headers, ",")))
		return nil, nil, errors.New(errmessage)

	} else {
		Logger.Infof("arquivo %s com formato correto para o processo", ExcelFile)
		rows = rows[1:]

		for _, row := range rows {

			pbi.Titulo = strings.TrimSpace(fmt.Sprintf("%s - %s", row[0], row[2]))
			pbi.Tags = row[1]
			pbi.Data = row[2]
			pbi.Name = row[3]
			pbi.Description = row[4]
			pbis = append(pbis, pbi)
		}
	}

	return pbis, f, nil
}

func checkHeadersExists(row []string, cabeçalhosEsperados []string) bool {

	// Compara os cabeçalhos obtidos com os esperados
	if len(row) != len(cabeçalhosEsperados) {
		Logger.Errorf("arquivo com número de informações acima do esperado tem %d e o esperado é %d", len(row), len(cabeçalhosEsperados))
		Logger.Errorf("verificar ordem dos headers e nomeclatura para seguir com o processo, headers esperados: %v", strings.ToLower(strings.Join(cabeçalhosEsperados, ",")))
		Logger.Errorf("headers no excel: %v", strings.ToLower(strings.Join(row, ",")))
		return false
	}
	for i, valor := range row {
		if !strings.EqualFold(valor, cabeçalhosEsperados[i]) {
			return false
		}
	}

	return true
}
func FinalizeAndSaveExcel(f *excelize.File, WorkItemCreatedUrl []*workitemtracking.WorkItem) error {

	// Insere uma nova coluna na primeira folha (Sheet1)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	f.SetCellValue(sheetName, "F1", "Result Sucess")
	for _, url := range WorkItemCreatedUrl {
		for i, row := range rows {

			if fmt.Sprintf("%s - %s", row[0], row[2]) == (*url.Fields)["System.Title"].(string) {
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", i), *url.Url)

				break
			}
		}

	}

	Logger.Debug("tasks criadas com sucesso ", WorkItemCreatedUrl)

	err = f.SaveAs("result.xlsx")
	if err != nil {
		Logger.Errorf("Error copying data: %v", err)
		return err
	}

	return nil
}

func CreateExcel(ReposInfoList []types.ReposInfo, Logger *logrus.Logger) {
	// cria um novo arquivo Excel
	f = excelize.NewFile()

	// cria uma nova planilha com o nome "azureDevops info"
	_, err := f.NewSheet("azureDevops info")
	if err != nil {
		Logger.Error(err)
	}
	headers := []string{"PROJECT", "DESCRIPTION", "VARIABLEGROUP", "VARIABLES"}
	for i, header := range headers {
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header)

	}

	for i, repoInfo := range ReposInfoList {
		for index, variableGroup := range repoInfo.VariableGroups {
			Logger.Debug(*variableGroup.Name)
			variables := mapToStringSlice(*variableGroup.Variables)
			in := i + 2 + index
			setCelValue("Sheet1", fmt.Sprintf("A%d", in), repoInfo.Name)
			setCelValue("Sheet1", fmt.Sprintf("B%d", in), repoInfo.RepoDescription)
			setCelValue("Sheet1", fmt.Sprintf("C%d", in), *variableGroup.Name)
			setCelValue("Sheet1", fmt.Sprintf("D%d", in), strings.Join(variables, ","))
		}

	}
	// escreve alguns dados nas células

	// salva o arquivo
	if err := f.SaveAs("azureDevopsInfo.xlsx"); err != nil {
		fmt.Println(err)
	}
}
func setCelValue(sheetname, cell, value string) error {
	err := f.SetCellValue(sheetname, cell, value)
	if err != nil {
		return err
	}
	return nil
}
func mapToStringSlice(m map[string]interface{}) []string {
	result := []string{}
	for _, value := range m {
		if str, ok := value.(string); ok {
			result = append(result, str)
		}
	}
	return result
}
