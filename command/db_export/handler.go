package db_export

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/filex/excelx"
	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"

	"quanx_tools/common/model"
)

func Export(dbType, dbName string, fields []*model.TableField) error {
	var tableMap = make(map[string][]*model.TableField)
	for _, field := range fields {
		tableMap[field.Table] = append(tableMap[field.Table], field)
	}
	var xlsxFile = xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return errorx.Wrap(err, "创建excel文件失败")
	}
	headers := excelx.GetHeadersByReflect(model.TableField{})
	var headerStyle = excelx.HeaderStyle()
	var defaultStyle = excelx.DefaultStyle()
	var headerCells []*xlsx.Cell
	for _, header := range headers {
		cell := &xlsx.Cell{Value: header.Name}
		cell.SetStyle(headerStyle)
		headerCells = append(headerCells, cell)
	}

	hMerge := len(headers) - 1
	for tableName, tableFields := range tableMap {
		// 添加表注释行
		tableCommentRow := sheet.AddRow()
		tableCommentCell := tableCommentRow.AddCell()
		tableCommentCell.Value = tableFields[0].TableComment
		tableCommentCell.HMerge = hMerge
		tableCommentCell.SetStyle(headerStyle)
		for i := 0; i < hMerge; i++ {
			tableCommentRow.AddCell().SetStyle(headerStyle)
		}
		// 添加表名行
		tableRow := sheet.AddRow()
		tableCell := tableRow.AddCell()
		tableCell.Value = tableName
		tableCell.HMerge = hMerge
		tableCell.SetStyle(headerStyle)
		for i := 0; i < hMerge; i++ {
			tableRow.AddCell().SetStyle(headerStyle)
		}

		// 添加表头行
		sheet.AddRow().Cells = headerCells
		for _, field := range tableFields {
			b, _ := json.Marshal(field)
			data := gjson.ParseBytes(b).Map()
			row := sheet.AddRow()
			for _, header := range headers {
				cell := row.AddCell()
				cell.SetStyle(defaultStyle)
				cell.Value = data[header.Key].String()
			}
		}
		// 添加间隔行
		sheet.AddRow()
	}
	// 这里重新生成
	dir := filepath.Join(dbType, dbName)
	filex.CreateDir(dir)
	if err = xlsxFile.Save(filepath.Join(dir, dbName+"_export_"+time.Now().Format("20060102150405")+".xlsx")); err != nil {
		return errorx.Wrap(err, "保存excel文件失败")
	}
	return nil
}
