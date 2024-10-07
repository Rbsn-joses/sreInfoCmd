package excel_test

import (
	"fmt"
	"testing"

	"github.com/Rbsn-joses/create-pbi/customLogger"
	"github.com/Rbsn-joses/create-pbi/excel"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

var err error
var logger = customLogger.InitLogger("info")

// Mock para excelize.OpenFile
func MockOpenFile(fileName string) (*excelize.File, error) {
	// Simule a abertura do arquivo com dados específicos
	if fileName == "correct_file.xlsx" {
		// Dados corretos
		return &excelize.File{}, nil
	} else if fileName == "incorrect_file.xlsx" {
		// Dados incorretos (número de colunas)
		return &excelize.File{}, nil
	} else if fileName == "incorrect_headers.xlsx" {
		// Dados incorretos (nomes de cabeçalhos)
		return &excelize.File{}, nil
	} else {
		return nil, fmt.Errorf("arquivo não encontrado: %s", fileName)
	}
}

// Mock para excelize.GetRows
func MockGetRows(sheetName string) ([]string, error) {
	if sheetName == "relatorio-pbi" {
		if sheetName == "correct_file.xlsx" {
			// Dados corretos (primeira linha com cabeçalhos esperados)
			return []string{"titulo", "tags", "data", "pbi", "description"}, nil
		} else if sheetName == "incorrect_headers.xlsx" {
			// Dados incorretos (nomes de cabeçalhos)
			return []string{"título", "tag", "data", "código pbi", "descrição"}, nil
		} else {
			// Dados incorretos (número de colunas)
			return []string{"titulo", "tags", "data"}, nil
		}
	} else {
		return nil, fmt.Errorf("planilha não encontrada: %s", sheetName)
	}
}
func TestGetExcel_HappyPath(t *testing.T) {
	fileName := "../correct_file.xlsx"
	_, _, err := excel.GetExcel(logger, fileName)

	assert.NoError(t, err)
}

func TestGetExcel_FileNotFound(t *testing.T) {
	fileName := "file-not-exist.xlsx"
	_, _, err := excel.GetExcel(logger, fileName)

	assert.Error(t, err)
	assert.EqualError(t, err, "open correct_file.xlsx: no such file or directory")
}

func TestGetExcel_IncorrectNumberOfColumns(t *testing.T) {
	fileName := "incorrect-number-headers.xlsx"
	pbis, _, err := excel.GetExcel(logger, fileName)

	assert.Contains(t, err, "verificar ordem dos headers e nomeclatura para seguir com o processo, headers esperados: titulo,tags,data,pbi,description")

	//assert.Error(t, err)
	assert.Nil(t, pbis)
}

func TestGetExcel_IncorrectHeaders(t *testing.T) {
	fileName := "incorrect-headers-sort.xlsx"
	_, _, err := excel.GetExcel(logger, fileName)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "headers esperados")
}

func TestGetExcel_SheetNotFound(t *testing.T) {
	fileName := "sheetNameIncorrect.xlsx"
	pbis, _, err := excel.GetExcel(nil, fileName)

	assert.Error(t, err)
	assert.Nil(t, pbis)
	assert.Contains(t, err.Error(), "planilha não encontrada")
}
