package model

import (
	"github.com/go-xuan/quanx/utils/randx"

	"quanx_tools/common"
)

// TableQuery 表基本信息
type TableQuery struct {
	Database string `json:"database"` // 数据库名
	Schema   string `json:"schema"`   // schema
	Table    string `json:"table"`    // 表名
	Columns  string `json:"columns"`  // 字段名
}

// TableDetail 表明细
type TableDetail struct {
	Database  string        `json:"database"`  // 数据库名
	Schema    string        `json:"schema"`    // schema
	Table     string        `json:"table"`     // 表名
	Columns   string        `json:"columns"`   // 字段名
	FieldList []*TableField `json:"fieldList"` // 字段列表
}

// TableField 表字段信息
type TableField struct {
	Name         string `json:"name" excel:"字段名"`      // 字段名
	Type         string `json:"type" excel:"数据类型"`     // 数据类型
	Precision    int    `json:"precision" excel:"长度"`  // 长度
	Scale        int    `json:"scale" excel:"小数点"`     // 小数点
	Nullable     bool   `json:"nullable" excel:"允许为空"` // 允许为空
	Default      string `json:"default" excel:"默认值"`   // 默认值
	Comment      string `json:"comment" excel:"注释"`    // 注释
	Database     string `json:"database"`              // 数据库名
	Schema       string `json:"schema"`                // schema名
	Table        string `json:"table"`                 // 表名
	TableComment string `json:"tableComment"`          // 表注释
}

type RandField struct {
	Name    string            `json:"name"`    // 字段名
	Type    string            `json:"type"`    // 数据类型
	Default string            `json:"default"` // 默认值
	Args    map[string]string `json:"args"`    // 约束参数
}

// RandField 将表字段结构体转为随机字段结构体
func (f *TableField) RandField() *RandField {
	return &RandField{
		Name:    f.Name,
		Type:    common.DB2RandType(f.Type),
		Default: f.Default,
	}
}

// RandFieldList 随机数生成字段配置
type RandFieldList []*RandField

func (list RandFieldList) RandResult(enums map[string][]string) map[string]string {
	result := make(map[string]string)
	for index, field := range list {
		var value string
		if field.Default == "" {
			var randModel = &randx.Options{
				Type:    field.Type,
				Param:   randx.NewParam(field.Args),
				Default: field.Default,
				Offset:  index,
				Enums:   enums[field.Name],
			}
			value = randModel.NewString()
		} else {
			value = field.Default
		}
		result[field.Name] = value
	}
	return result
}
