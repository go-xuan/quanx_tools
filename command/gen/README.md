# 1.配置文件

#### 配置文件说明

| 配置项         | 必填    | 说明                            |
|:------------|-------|-------------------------------|
| app         | true  | 应用名                           |
| frame       | true  | 项目框架名, 自定义配置                  |
| template    | false | 外部自定义模板路径，默认“./template”      |
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

#### 配置文件示例

```yaml
# 应用名
app: demo
# 模板路径
template:
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
# 2.模板文件

## 2.1 模板文件规则
- 框架模板分为“内置框架模板”和“外置框架模板”两类。代码生成工具仅内置了【go-quanx】框架模版，可以直接使用。在工具执行时，如果“内置框架模板”和“外置框架模板”同名，则优先基于“外置框架模板”生成代码。
- 如需使用“外置框架模板”，需要在config.yaml中配置模板根文件夹，即template配置项，默认值“./template”。
- 在模板根文件夹下，可以配置不限数量的“框架模板”，文件夹名即框架名，支持自定义，例如：“./template/aaa/”文件夹即aaa框架，在此目录下新增其代码模板文件即可。
- “框架模板”即一个文件夹，文件夹下包含不限数量的模板文件。模板文件可以是任何文件名或者任意文件类型，但是必须使用".tmpl"作为文件后缀，例如，想要生成main.go代码，模板文件名需要命名为main.go.tmpl；想要生成hello.java代码，模板文件名需要命名为hello.java.tmpl。
- 代码生成工具完全按照框架模板内的目录结构进行代码生成，即生成代码的目录结构与框架模板内的目录结构保持一致， 例如：模板文件template/{{frame}}/config/config.yaml.tmpl生成的代码为：{{output}}/{{app}}/config/config.yaml。
- 代码生成工具会根据模板文件的文件名是否含有{{model}}占位符，执行不同的代码生成策略。如果模板文件名携带了{{model}}占位符，则会基于当前所有的表结构模型生成一份代码文件，例：假设模板文件为{{model}}.go.tmpl，如果本次执行gen表模型包含aaa、bbb、ccc三张表，则会对应生成aaa.go、bbb.go、ccc.go代码。
- 模板文件第一行加上“// this file will be overwritten when execute gen command next.”注释，下一次执行代码生成工具时，当识别到此行注释则会重新生成当前代码文件进行覆盖。

## 2.2 模板目录结构

- 以下模板目录结构仅为参考示例
```shell
└── template                                    # 模板根文件夹	
       ├── frame_aaa                            # 框架模板aaa
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
       ├── frame_bbb                            # 框架模板bbb
       │     ├── ....                           # 框架模板bbb的模板文件
       │     ├── ....                           
       ├── ......                               # 其他自定义框架模板 
       └── frame_xxx                               
```

## 2.3 模板占位符变量

### 2.3.1 模板文件名不带{{model]}

```gotemplate
# 以下当前模板文件可支持的占位符变量以及变量含义
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

# 对所有的表结构模型{{.Models}} 进行遍历
 {{- range $model := .Models}}
    {{$model.Table}}          # 数据库原始表名
    {{$model.Name}}           # 模型名（原始表名去除前后缀）
    {{$model.Schema}}         # 数据库模式名（pg专有）
    {{$model.Comment}}        # 表注释
    {{$model.Fields}}         # 所有表字段，需要遍历后再进行取值
   
    # 对当前表的字段{{.Fields}} 进行遍历
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
### 2.3.2 模板文件名带{{model]}

- 代码生成工具识别到模板文件的文件名包含{{model}}占位符，即认为当前模板文件是基于单表模型生成代码，自动遍历当前所有的表模型，为每个表模型都生成一份代码。 

```gotemplate
# 以下当前模板文件可支持的占位符变量以及变量含义
{{.App}}            # config.yaml配置的app值
{{.Table}}          # 数据库表名 
{{.Name}}           # 模型名（数据库去除前后缀）
{{.Schema}}         # 数据库模式名（pg专有）     
{{.Comment}}        # 表注释
{{.Fields}}         # 当前表字段，需要遍历后再进行取值

# 对当前表模型{{model}}}的字段{{.Fields}} 进行遍历
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

## 2.4 字符转化函数

### 2.4.1 使用规则

```shell
# 目前支持函数
uc        # 转大驼峰，例：table_name  => TableName
lc        # 转小驼峰，例：table_name  => tableName
snake     # 转蛇形，例：TableName  =>  table_name
path      # 转路径，例：table_name  => table/name
```

#### 2.4.2 使用示例

```gotemplate
// 示例
router.POST("/api/v1/{{path .Name}}", {{uc .Name}}Handler)
```

## 2.5 模板文件覆盖

#### 2.5.1 使用规则

- 对于和表结构强关联的代码，可在模板文件首行添加“// this file will be overwritten when execute gen command next.”注释，下次代码生成工具重新执行的时候，则会对该代码文件内容进行覆盖重写。
- 如果想代码文件在下次执行中不受重写影响，则需要移除此行注释


#### 2.5.2 使用示例

```gotemplate
// this file will be overwritten when execute gen command next.

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