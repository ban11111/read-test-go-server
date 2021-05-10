package exporter

import (
	"fmt"
	"github.com/dipperin/go-ms-toolkit/json"
	"github.com/stretchr/testify/assert"
	"read-test-server/model"
	"testing"
)

func TestExcelExporter_makeContent(t *testing.T) {
	e := &ExcelExporter{}

	datas := []*model.User{{Name: "1"}, {Name: "2"}}

	content, err := e.makeContent(datas)

	assert.NoError(t, err)
	fmt.Println(json.StringifyJson(content))
}

func TestExcelExporter_makeTitle(t *testing.T) {
	e := &ExcelExporter{}
	datas := []*model.User{{Name: "1"}, {Name: "2"}}

	title := e.makeTitle(datas)
	fmt.Println(json.StringifyJson(title))
}