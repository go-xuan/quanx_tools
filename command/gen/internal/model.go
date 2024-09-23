package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-xuan/quanx/types/stringx"
)

type Model struct {
	App          string   `json:"app"`          // 应用名
	Table        string   `json:"table"`        // 表名
	Name         string   `json:"name"`         // 模型名
	Schema       string   `json:"schema"`       // schema
	Comment      string   `json:"comment"`      // 表备注
	FiledNameLen int      `json:"filedNameLen"` // 字段名称长度
	FiledTypeLen int      `json:"filedTypeLen"` // 字段类型长度
	HasTime      bool     `json:"hasTime"`      // 是否含有时间类型
	Fields       []*Field `json:"fields"`       // 字段列表
}

type Field struct {
	Name         string `json:"name"`         // 字段名
	Type         string `json:"type"`         // 数据类型
	Precision    int    `json:"precision"`    // 长度
	Scale        int    `json:"scale"`        // 小数点
	Nullable     bool   `json:"nullable"`     // 允许为空
	Default      string `json:"default"`      // 默认值
	Comment      string `json:"comment"`      // 注释
	Database     string `json:"database"`     // 数据库名
	Schema       string `json:"schema"`       // schema名
	TableName    string `json:"tableName"`    // 表名
	TableComment string `json:"tableComment"` // 表注释
	GoName       string `json:"goName"`       // go字段名
	GoType       string `json:"goType"`       // go数据类型
	GormType     string `json:"ormType"`      // orm字段类型
	JavaType     string `json:"javaType"`     // Java字段类型
}

func (t *Model) SelectSql() string {
	ss := "%" + strconv.Itoa(-t.FiledNameLen) + "s as %s,"
	sb := strings.Builder{}
	sb.WriteString("select ")
	for i, field := range t.Fields {
		if i > 0 {
			sb.WriteString("\n")
			sb.WriteString("       ")
		}
		sb.WriteString(fmt.Sprintf(ss, field.Name, stringx.ToLowerCamel(field.Name)))
	}
	sb.WriteString(",\n  from ")
	sb.WriteString(t.Name)
	return strings.ReplaceAll(sb.String(), ",,", "")
}

func (t *Model) InsertSql() string {
	sb, iv := strings.Builder{}, strings.Builder{}
	sb.WriteString("insert ")
	sb.WriteString("into ")
	sb.WriteString(t.Name)
	sb.WriteString("\n(")
	var i int
	for _, field := range t.Fields {
		var fieldName = field.Name
		if ignoreField(fieldName) && field.Default != "" {
			continue
		}
		if i > 0 {
			sb.WriteString("\n")
			iv.WriteString("\n")
		}
		sb.WriteString(fieldName)
		sb.WriteString(",")
		iv.WriteString("#{create.")
		iv.WriteString(stringx.ToLowerCamel(fieldName))
		iv.WriteString("},")
		i++
	}
	sb.WriteString(",)\nvalues \n(")
	sb.WriteString(iv.String())
	sb.WriteString(",)")
	return strings.ReplaceAll(sb.String(), ",,", "")
}

func (t *Model) UpdateSql() string {
	sb := strings.Builder{}
	sb.WriteString("update ")
	sb.WriteString(t.Name)
	sb.WriteString("\n<set>")
	for _, field := range t.Fields {
		var fieldName = field.Name
		if ignoreField(fieldName) {
			continue
		}
		lc := "update." + stringx.ToLowerCamel(fieldName)
		sb.WriteString("\n\t")
		if field.Type == Varchar {
			sb.WriteString(fmt.Sprintf(`<if test="%s != null and %s != ''"> %s = #{%s}, </if>`, lc, lc, fieldName, lc))
		} else {
			sb.WriteString(fmt.Sprintf(`<if test="%s != null"> %s = #{%s}, </if>`, lc, fieldName, lc))
		}
	}
	sb.WriteString("\n</set>")
	sb.WriteString("\nwhere id = #{update.id}")
	return strings.ReplaceAll(sb.String(), ",,", "")
}

// BASE基础字段映射
func ignoreField(field string) bool {
	switch field {
	case "id":
		return true
	case "create_time":
		return true
	case "update_time":
		return true
	case "create_user_id":
		return true
	case "update_user_id":
		return true
	case "create_by":
		return true
	case "update_by":
		return true
	default:
		return false
	}
}
