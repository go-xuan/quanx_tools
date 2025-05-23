package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/extra/gormx"
	"github.com/go-xuan/quanx/types/stringx"

	embedTemplate "quanx_tools/command/gen/template"
)

type Config struct {
	App      string   `json:"app" json:"app" default:"demo"`         // 应用名
	Frame    string   `json:"frame" yaml:"frame" default:"go-quanx"` // 模板框架（go-quanx/go-frame/spring-boot）
	Template string   `json:"template" yaml:"template"`              // 外置模板路径
	Output   string   `json:"output" yaml:"output"`                  // 输出路径
	DB       DBConfig `json:"db" yaml:"db"`                          // 数据库
}

type DBConfig struct {
	Type        string `json:"type" yaml:"type"`                     // 数据库类型
	Host        string `json:"host" yaml:"host" default:"localhost"` // 数据库Host
	Port        int    `json:"port" yaml:"port"`                     // 数据库端口
	Username    string `json:"username" yaml:"username"`             // 用户名
	Password    string `json:"password" yaml:"password"`             // 密码
	Database    string `json:"database" yaml:"database"`             // 数据库名
	Schema      string `json:"schema" yaml:"schema"`                 // schema模式名
	Include     string `json:"include" yaml:"include"`               // 包含表（为空则获取全表,多个以“,”分割）
	Exclude     string `json:"exclude" yaml:"exclude"`               // 排除表（多个以“,”分割）
	TablePrefix string `json:"tablePrefix" yaml:"tablePrefix"`       // 表前缀
	TableSuffix string `json:"tableSuffix" yaml:"tableSuffix"`       // 表后缀缀
}

func (db DBConfig) GormxDB() *gormx.Config {
	if db.Type != "" && db.Host != "" {
		// 初始化数据库连接
		gormDB := &gormx.Config{
			Enable:   true,
			Type:     db.Type,
			Host:     db.Host,
			Port:     db.Port,
			Username: db.Username,
			Password: db.Password,
			Database: db.Database,
			Schema:   db.Schema,
		}
		if err := gormDB.Execute(); err != nil {
			fmt.Println("连接数据库失败：", err)
			return nil
		}
		return gormDB
	}
	return nil
}

func (c *Config) Root() string {
	if strings.HasSuffix(c.Output, c.App) {
		return c.Output
	} else {
		return filepath.Join(c.Output, c.App)
	}
}

// ExternalTemplateCheck 外置模板检测
func (c *Config) ExternalTemplateCheck() error {
	fmt.Println("外置模板检测框架：", c.Frame)
	dir := filepath.Join(stringx.IfZero(c.Template, TemplateDir), c.Frame)
	fmtx.Green.Xprintf("外置模板检测路径：%s\n", dir)
	if filex.Exists(dir) {
		if files, err := filex.FileScan(dir, filex.OnlyFile); err == nil {
			for _, file := range files {
				path := file.Path
				fmtx.Green.Xprintf("当前检测模板文件：%s\n", path)
				if filex.GetSuffix(path, true) != embedTemplate.Suffix {
					fmtx.Red.Xprintf("模板文件后缀异常：%s\n", path)
					if err = os.Rename(path, path+embedTemplate.Suffix); err != nil {
						return errorx.Wrap(err, "模板文件后缀矫正失败："+path)
					} else {
						fmtx.Red.Println("模板文件后缀矫正成功")
					}
				}
			}
		} else {
			return errorx.Wrap(err, "模板扫描失败："+dir)
		}
	} else {
		fmtx.Red.Println("未发现模板文件，请检查后重试！")
	}
	return nil
}

func (c *Config) Generator() *Generator {
	var generator = &Generator{App: c.App, Root: c.Root()}
	// 获取代码生成模板文件
	var templateDir = stringx.IfZero(c.Template, TemplateDir)
	if files := GetExternalTemplateFiles(templateDir, c.Frame); len(files) > 0 {
		fmt.Println("读取外置模板：", templateDir, c.Frame)
		generator.TmplFiles = files
	} else if files = GetInternalTemplateFiles(embedTemplate.FS, c.Frame, c.Frame); len(files) > 0 {
		fmt.Println("读取内置模板：", c.Frame)
		generator.TmplFiles = files
	}
	// 查询数据库表模型失败
	if gormDB := c.DB.GormxDB(); gormDB != nil {
		generator.DB = *gormDB
		if models, err := c.DB.GetModels(c.App); err == nil {
			generator.Models = models
		} else {
			fmtx.Red.Xprintf("查询数据库表模型失败: %s \n", err.Error())
		}
	}
	return generator
}

// Trim 去除表名前后缀
func (db DBConfig) Trim(table string) string {
	if db.TablePrefix != "" {
		table = strings.TrimPrefix(table, db.TablePrefix)
	}
	if db.TableSuffix != "" {
		table = strings.TrimSuffix(table, db.TableSuffix)
	}
	return table
}

// GetModels 查询数据库表模型
func (db DBConfig) GetModels(app string) ([]*Model, error) {
	defer gormx.Close()
	var sql string
	switch db.Type {
	case gormx.MYSQL:
		sql = mysqlTableFieldQuerySql(db.Database, db.Include, db.Exclude)
	case gormx.POSTGRES, gormx.PGSQL:
		sql = pgsqlTableFieldQuerySql(db.Database, db.Schema, db.Include, db.Exclude)
	default:
		return nil, errorx.New("数据库类型（db.type）只支持mysql或者postgres")
	}
	var fields []*Field
	if err := gormx.GetInstance().Raw(sql).Scan(&fields).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表字段列表失败")
	}
	if len(fields) > 0 {
		var modelMap = make(map[string]*Model)
		for _, field := range fields {
			if field.Default == "auto_increment" ||
				strings.Contains(field.Default, "::") {
				field.Default = ""
			}
			field.GoName = stringx.ToUpperCamel(field.Name)
			field.GoType = DB2GoType(field.Type)
			field.GormType = DB2GormType(field.Type)
			field.JavaType = DB2JavaType(field.Type)
			var nameLen, typeLen = len(field.GoName), len(field.GoType)
			var table = field.TableName
			if _, ok := modelMap[table]; ok {
				modelMap[table].Fields = append(modelMap[table].Fields, field)
				if modelMap[table].FiledNameLen < nameLen {
					modelMap[table].FiledNameLen = nameLen
				}
				if modelMap[table].FiledTypeLen < typeLen {
					modelMap[table].FiledTypeLen = typeLen
				}
			} else {
				modelMap[table] = &Model{
					App:          app,
					Table:        table,
					Name:         db.Trim(table),
					Database:     field.Database,
					Schema:       field.Schema,
					Comment:      field.TableComment,
					FiledNameLen: nameLen,
					FiledTypeLen: typeLen,
					Fields:       []*Field{field},
				}
			}
			if field.GoType == GoTime {
				modelMap[table].HasTime = true
			}
		}
		var models []*Model
		for _, model := range modelMap {
			for _, field := range model.Fields {
				field.GoName = stringx.Fill(field.GoName, " ", model.FiledNameLen)
				field.GoType = stringx.Fill(field.GoType, " ", model.FiledTypeLen)
			}
			models = append(models, model)
		}
		return models, nil
	}
	return nil, nil
}
