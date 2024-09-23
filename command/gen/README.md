# 配置文件

### 配置文件说明

| 配置项         | 必填    | 说明                            |
|:------------|-------|-------------------------------|
| app         | true  | 应用名                           |
| frame       | true  | 项目框架名, 自定义配置                  |
| templateDir | false | 外部自定义模板路径，默认“./template”      |
| output      | true  | 输出路径，即生成代码保存位置                |
| include     | false | 使用表，为空则获取全表，多个以“,”分割          |
| exclude     | false | 排除表，多个以“,”分割                  |
| tablePrefix | false | 表名前缀，生成model模型时，结构体名不会带上此前缀   |
| tableSuffix | false | 表名后缀，生成model模型时，结构体名不会带上此前缀   |
| db          | false | 数据库配置，如果需要根据现有表结构生成代码，可配合此数据源 |
| -type       | false | 数据库类型（mysql/postgres）         |
| -host       | false | 数据库host                       |
| -port       | false | 端口                            |
| -username   | false | 用户名                           |
| -password   | false | 密码                            |
| -database   | false | 数据库名                          |
| -schema     | false | 模式名（pg专用）                     |

### 配置文件示例
```yaml
# 应用名
app: demo
# 模板路径
templateDir:
# 项目框架
frame: go-quanx
# 输出路径
output: ./
# 包含表名（为空则获取全表，多个以“,”分割）
include:
# 排除表名（多个以“,”分割）
exclude:
# 去除表名前缀
tablePrefix: 
# 去除表名后缀
tableSuffix:
# 数据库配置
db:
  type: mysql
  host: localhost
  port: 3306
  username: root
  password: root
  database: demo
```
# 模板文件

## 模板文件规则
- 如需使用外部模板，需要在config.yaml中配置templateDir
- 在模板根路径下，可以配置不限数量的自定义明明的框架模板，文件夹名即frame框架名
- 模板文件可以是任何文件名或者任意文件类型，但是必须使用".tmpl"作为文件后缀，例：如需生成main.go代码，对应模板文件名则为main.go.tmpl；如需生成hello.java代码，对应模板文件名则为hello.java.tmpl
- 代码生成工具仅内置了go-quanx的框架模版，可直接使用。在工具执行时，如果外部模板根路径下存在同名框架模板，则会优先根据外部模板生成代码。
- 代码生成工具完全按照模板内自定义配置的目录结构进行代码生成，生成后的代码目录结构与模板目录保持一致， 例如：模板文件{{templateDir}}/{{frame}}/config/config.yaml.tmpl生成的代码为：{{output}}/{{app}}/config/config.yaml
- 代码生成工具会根据模板文件名是否包含{{model}}占位符，判断是否需要进行表结构遍历，以便基于表结构生成相应代码，例：数据库如果存在aaa、bbb、ccc三张表，则会对应生成aaa.go、bbb.go、ccc.go 代码
- 模板文件第一行加上“// This file will be overwritten on re-execution.”注释，下次执行代码生成工具时，识别到此行注释则会进行覆盖。

## 模板目录结构
- 以下模板目录结构仅为示例

```shell
└── template                                    # 模板根文件夹	
       ├── frame_a                               # 框架模板1
       │     ├── config                         # 配置文件目录
       │     │     ├── config.yaml.tmpl         
       │     │     ├── database.yaml.tmpl           
       │     │     ├── ......                   # 其他配置文件
       │     │     └── xxx.yaml.tmpl          
       │     ├── service                        # 业务代码目录 
       │     │     ├── service_a.go.tmpl        # 业务代码	
       │     │     ├── service_b.go.tmpl        
       │     │     ├── service_c.go.tmpl        
       │     │     ├── .....                    # 其他业务代码
       │     │     └── xxx.go.tmpl            	
       │     ├── model                          # 模型代码目录
       │     │     ├── common.go.tmpl           # 通用模型	
       │     │     └── {{model}}.go.tmpl        # 表结构模型，文件名需要带上{{model}}占位符，否则工具将会识别为普通模板
       │     ├── xxx.yyy.tmpl                   # 其他格式文件模板，xxx表示任意文件名，yyy表示任意文件后缀名，最终以".tmpl"作为后缀即可，否则工具无法识别为模板
       │     ├── .....                          # 其他文件		
       │     ├── go.mod.tmpl                    # go mod模块名	
       │     └── main.go.tmpl                   # main主程序入口	
       ├── frame_b                               
       ├── ......                               # 其他自定义框架模板 
       └── xxx                               
```

## 模板占位符

### 模板文件名不带{{model]}

```gotemplate
{{.App}}              # config.yaml配置的app值      
{{.DB}}               # config.yaml配置的数据库配置
{{.Models}}           # 所有表结构模型，需要遍历后再进行取值

# 对{{.DB}} 的内部元素进行取值   
{{.DB.Type}}          # 数据库类型（mysql/postgres）
{{.DB.Host}}          # 数据库Host
{{.DB.Port}}          # 数据库端口
{{.DB.Username}}      # 用户名
{{.DB.Password}}      # 密码
{{.DB.Database}}      # 数据库名
{{.DB.Schema}}        # schema模式名（pg专有）

# 对{{.Models}} 进行遍历之后，对表结构元素进行取值
 {{- range $model := .Models}}
    {{$model.Table}}          # 数据库原始表名
    {{$model.Name}}           # 模型名（原始表名去除前后缀）
    {{$model.Schema}}         # 数据库模式名（pg专有）
    {{$model.Comment}}        # 表注释
    {{$model.Fields}}         # 所有表字段，需要遍历后再进行取值
   
    # 对{{.Fields}} 进行遍历之后，对表字段的元素进行取值
    {{- range $Field := .Fields}}
        {{$Field.Name}}             # 字段名 
        {{$Field.Type}}             # 字段类型
        {{$Field.Precision}}        # 字段长度  
        {{$Field.Scale}}            # 字段精度 
        {{$Field.Nullable}}         # 允许为空
        {{$Field.Default}}          # 字段默认值 
        {{$Field.Comment}}          # 字段注释  
    {{- end}}
 {{- end}}
```
### 模板文件名带{{model]}
- 代码生成工具识别模板文件名包含{{model}}占位符，即认为当前模板文件在所在目录下，需要根据每张表结构都生成一份代码。 
```gotemplate
{{.App}}            # config.yaml配置的app值
{{.Table}}          # 数据库表名 
{{.Name}}           # 模型名（数据库去除前后缀）
{{.Schema}}         # 数据库模式名（pg专有）     
{{.Comment}}        # 表注释
{{.Fields}}         # 当前表字段，需要遍历后再进行取值

# 对{{.Fields}} 进行遍历之后，对表字段的元素进行取值
 {{- range $Field := .Fields}}
    {{$Field.Name}}             # 字段名 
    {{$Field.Type}}             # 字段类型
    {{$Field.Precision}}        # 字段长度  
    {{$Field.Scale}}            # 字段精度 
    {{$Field.Nullable}}         # 允许为空
    {{$Field.Default}}          # 字段默认值 
    {{$Field.Comment}}          # 字段注释  
 {{- end}}
```

## 字符转化函数

#### 支持函数

```shell
# 目前支持函数
uc        # 转大驼峰，例：table_name  => TableName
lc        # 转小驼峰，例：table_name  => tableName
snake     # 转蛇形，例：TableName  =>  table_name
path      # 转路径，例：table_name  => table/name
```

#### 使用方式

```gotemplate
// 示例
router.POST("/api/v1/{{path .Name}}", {{uc .Name}}Handler)
```

### 模板文件覆盖
- 对于和表结构强关联的代码，可在模板文件首行添加“// This file will be overwritten on re-execution.”注释，下次代码生成工具重新执行的时候则会对该文件内容进行重写。
- 如果想在下次执行中不受重写影响，则需要移除此行注释

```gotemplate
// This file will be overwritten on re-execution.

package sqlx_model

import (
	"strings"
	"time"
)

type {{uc .Name}} struct {
	{{- range $field := .Fields}}
	{{$field.GoName}} {{$field.GoType}} `db:"{{$field.Name}}" json:"{{$field.Name}}"` // {{$field.Comment}}
	{{- end }}
}

func (u *{{uc .Name}}) TableName() string {
	return "{{.Name}}"
}

```