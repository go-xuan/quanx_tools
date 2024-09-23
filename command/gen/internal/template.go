package internal

import (
	"bufio"
	"bytes"
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/fmtx"

	embedTemplate "quanx_tools/command/gen/template"
)

var funcs = template.FuncMap{
	"uc":    stringx.ToUpperCamel,
	"lc":    stringx.ToLowerCamel,
	"snake": stringx.ToSnake,
	"path":  ToPath,
}

// ToPath 转路径
func ToPath(s string) string {
	s = stringx.ToSnake(s)
	s = strings.ReplaceAll(s, "_", "/")
	return s
}

// TemplateFile 代码生成模板
type TemplateFile struct {
	Frame    string           // 框架名
	Path     string           // 模版路径
	Content  string           // 模板内容
	DataType string           // 填充数据
	FuncMap  template.FuncMap // 自定义模板参数方法
}

// WriteDataToFile 写入数据到模板文件中
func (t *TemplateFile) WriteDataToFile(root string, data any, model ...string) error {
	var dir, file = filepath.Split(t.Path)
	file = strings.TrimSuffix(file, embedTemplate.Suffix)
	dir = strings.TrimPrefix(dir, t.Frame)
	if len(model) > 0 {
		file = strings.Replace(file, "{{model}}", model[0], -1)
	}
	filePath := filepath.Join(root, dir, file)
	if writeMode := doWriteOrNot(filePath); writeMode != dontWrite {
		var buf = &bytes.Buffer{}
		tt := template.Must(template.New(t.Path).Funcs(t.FuncMap).Parse(t.Content))
		if err := tt.Execute(buf, data); err != nil {
			return errorx.Wrap(err, "模版执行失败"+t.Path)
		}
		if err := filex.WriteFile(filePath, buf.String()); err != nil {
			return errorx.Wrap(err, "写入文件失败"+filePath)
		}
		if writeMode == doOverwrite {
			fmtx.Green.XPrintf("代码覆盖成功：%s \n", filePath)
		} else {
			fmtx.Green.XPrintf("代码生成成功：%s \n", filePath)
		}
	}
	return nil
}

// CustomTemplateFiles 获取自定义模板
func CustomTemplateFiles(dir, frame string) []*TemplateFile {
	//先从外部template文件夹
	if fDir := filepath.Join(dir, frame); filex.Exists(fDir) {
		if files, err := filex.FileScan(fDir, filex.OnlyFile); err == nil {
			for _, file := range files {
				dataType := GeneratorData
				if strings.Contains(file.Info.Name(), "{{model}}") {
					dataType = ModelData
				}
				var content, _ = filex.ReadFile(file.Path)
				path := strings.TrimPrefix(file.Path, dir+string(os.PathSeparator))
				var templates []*TemplateFile
				templates = append(templates, &TemplateFile{
					Frame:    frame,
					Path:     path,
					Content:  string(content),
					DataType: dataType,
					FuncMap:  funcs,
				})
				return templates
			}
		}
	}
	return nil
}

// EmbedTemplateFiles 获取内置模板
func EmbedTemplateFiles(fs embed.FS, dir, frame string) []*TemplateFile {
	var templates []*TemplateFile
	if entries, err := fs.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				var entryDir = filepath.Join(dir, entry.Name())
				templates = append(templates, EmbedTemplateFiles(fs, entryDir, frame)...)
			} else {
				var fileName, dataType = entry.Name(), GeneratorData
				if strings.Contains(fileName, "{{model}}") {
					dataType = ModelData
				}
				var path = filepath.Join(dir, fileName)
				var content, _ = fs.ReadFile(path)
				templates = append(templates, &TemplateFile{
					Frame:    frame,
					Path:     path,
					Content:  string(content),
					DataType: dataType,
					FuncMap:  funcs,
				})
			}
		}
	}
	return templates
}

const (
	dontWrite = iota
	doWrite
	doOverwrite
)

// 是否跳过文件写入
func doWriteOrNot(path string) int {
	if filex.Exists(path) {
		var err error
		var file *os.File
		if file, err = os.OpenFile(path, os.O_RDONLY, 0666); err != nil {
			return doWrite
		}
		defer file.Close()
		var line []byte
		if line, _, err = bufio.NewReader(file).ReadLine(); err != nil {
			return doWrite
		} else {
			if string(line) == OverwriteTag {
				return doOverwrite
			} else {
				return dontWrite
			}
		}
	}
	return doWrite
}
