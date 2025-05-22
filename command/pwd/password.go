package pwd

import (
	"embed"
	"fmt"
	"io/fs"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
	"quanx_tools/common/table"
)

//go:embed *
var FS embed.FS

var (
	Command = flagx.NewCommand(command.Password, "密码")
	data    = enumx.NewEnum[string, []*Password]()
)

func (p *Password) ToString() string {
	return fmtx.Yellow.Xsprintf(
		"username: %s password: %s database: %s port: %s  host: %s",
		stringx.Fill(p.Username, " ", 30),
		stringx.Fill(p.Password, " ", 30),
		stringx.Fill(p.Database, " ", 30),
		strconv.Itoa(p.Port),
		p.Host,
	)
}

func init() {
	Command.AddOption(
		flagx.BoolOption("json", "json格式化", false),
		flagx.BoolOption("yaml", "yaml格式化", false),
	)

	_ = fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filex.GetSuffix(path) != "go" {
			var content []byte
			if content, err = FS.ReadFile(path); err != nil {
				return err
			}
			var pwds []*Password
			if err = marshalx.Apply(path).Unmarshal(content, &pwds); err != nil {
				return err
			}
			_, filename, _ := filex.Analyse(path)
			data.Add(filename, pwds)
			Command.AddOption(flagx.BoolOption(filename, fmt.Sprintf("获取%s密码", filename), false))
		}
		return nil
	})

	Command.AddOption(flagx.StringOption("keyword", "检索关键字", ""))
	Command.SetExecutor(executor)
}

func executor() error {
	var pwds []*Password
	for _, name := range data.Keys() {
		if Command.GetOptionValue(name).Bool() {
			pwds = data.Get(name)
		}
	}
	if pwds == nil || len(pwds) == 0 {
		Command.OptionsHelp()
		return nil
	}

	// 获取指定数据源
	var output []*Password
	if keyword := Command.GetOptionValue("keyword").String(); keyword != "" {
		for _, pwd := range pwds {
			if strings.Contains(pwd.Source, keyword) || strings.Contains(pwd.Database, "keyword") {
				output = append(output, pwd)
			}
		}
	} else {
		output = pwds
	}

	// 输出结果
	if Command.GetOptionValue("json").Bool() {
		bytes, _ := marshalx.Json("    ").Marshal(output)
		fmt.Println(string(bytes))
	} else if Command.GetOptionValue("yaml").Bool() {
		bytes, _ := marshalx.Yaml().Marshal(output)
		fmt.Println(string(bytes))
	} else {
		table.Print(output)
	}
	return nil
}

type Password struct {
	Source   string `yaml:"source" json:"source" column:"source"`       // 数据源
	Type     string `yaml:"type" json:"type"`                           // 数据库类型
	Host     string `yaml:"host" json:"host" column:"host"`             // host
	Port     int    `yaml:"port" json:"port" column:"port"`             // 端口
	Database string `yaml:"database" json:"database" column:"database"` // 数据库
	Username string `yaml:"username" json:"username" column:"username"` // 用户名
	Password string `yaml:"password" json:"password" column:"password"` // 密码
}

func (p *Password) GetHeaders() []string {
	valueOf := reflect.ValueOf(p).Elem()
	typeOf := valueOf.Type()
	var names []string
	for i := 0; i < typeOf.NumField(); i++ {
		if name, ok := typeOf.Field(i).Tag.Lookup("column"); ok {
			names = append(names, name)
		}
	}
	return names
}

func (p *Password) GetValuesAndWides() ([]string, []int) {
	valueOf := reflect.ValueOf(p).Elem()
	typeOf := valueOf.Type()
	var values []string
	var wides []int
	for i := 0; i < typeOf.NumField(); i++ {
		if _, ok := typeOf.Field(i).Tag.Lookup("column"); ok {
			value := fmt.Sprintf("%v", valueOf.Field(i).Interface())
			values = append(values, value)
			wides = append(wides, len(value))
		}
	}
	return values, wides
}
