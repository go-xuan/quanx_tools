package internal

import (
	"github.com/go-xuan/quanx/core/gormx"
	"github.com/go-xuan/quanx/os/errorx"
)

// Generator 代码生成器
type Generator struct {
	App       string          // 应用名
	Root      string          // 代码生成路径
	DB        gormx.Config    // 应用数据库
	TmplFiles []*TemplateFile // 模板
	Models    []*Model        // 模型列表
}

// Execute 生成代码
func (gen *Generator) Execute() error {
	if len(gen.TmplFiles) == 0 {
		return errorx.Errorf("模板文件为空")
	}
	for _, file := range gen.TmplFiles {
		switch file.DataType {
		case EmptyData:
			if err := file.WriteDataToFile(gen.Root, nil); err != nil {
				return errorx.Wrap(err, "根据模板生成代码失败")
			}
		case ModelData:
			for _, model := range gen.Models {
				if err := file.WriteDataToFile(gen.Root, model, model.Name); err != nil {
					return errorx.Wrap(err, "根据模板生成代码失败")
				}
			}
		case GeneratorData:
			if err := file.WriteDataToFile(gen.Root, gen); err != nil {
				return errorx.Wrap(err, "根据模板生成代码失败")
			}
		}
	}
	return nil
}
