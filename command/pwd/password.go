package pwd

import (
	"embed"
	"fmt"
	"io/fs"
	"strconv"

	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
)

//go:embed *
var FS embed.FS

var (
	Command = flagx.NewCommand(command.Password, "密码")
	names   = enumx.NewEnum[string, string]()
)

type Password struct {
	Source   string `yaml:"source" json:"source"`     // 数据源
	Type     string `yaml:"type" json:"type"`         // 数据库类型
	Host     string `yaml:"host" json:"host"`         // host
	Port     int    `yaml:"port" json:"port"`         // 端口
	Database string `yaml:"database" json:"database"` // 数据库
	Username string `yaml:"username" json:"username"` // 用户名
	Password string `yaml:"password" json:"password"` // 密码
}

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
			_, name, _ := filex.Analyse(path)
			Command.AddOption(flagx.BoolOption(name, fmt.Sprintf("获取%s密码", name), false))
			names.Add(name, path)
		}
		return nil
	})

	Command.AddOption(flagx.StringOption("source", "选择数据源", ""))
	Command.SetExecutor(executor)
}

func GetPwdEnumFromFS(path string) (*enumx.Enum[string, *Password], error) {
	content, err := FS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pwds []*Password
	if err = marshalx.Apply(path).Unmarshal(content, &pwds); err != nil {
		return nil, err
	}

	var enum = enumx.NewEnum[string, *Password]()
	for _, pwd := range pwds {
		enum.Add(pwd.Source, pwd)
	}
	return enum, nil
}

func executor() error {
	var path string
	for _, name := range names.Keys() {
		if Command.GetOptionValue(name).Bool() {
			path = names.Get(name)
		}
	}
	if path == "" {
		Command.OptionsHelp()
		return nil
	}

	// 从FS获取所有账密
	pwds, err := GetPwdEnumFromFS(path)
	if err != nil {
		return err
	}

	// 序列化方式
	var marshal marshalx.Method
	if Command.GetOptionValue("json").Bool() {
		marshal = marshalx.Json("    ")
	} else if Command.GetOptionValue("yaml").Bool() {
		marshal = marshalx.Yaml()
	}

	if pwd := pwds.Get(Command.GetOptionValue("source").String()); pwd != nil {
		if marshal != nil {
			bytes, _ := marshal.Marshal(pwd)
			fmt.Println(string(bytes))
		} else {
			fmt.Println(pwd.ToString())
		}
	} else {
		if marshal != nil {
			bytes, _ := marshal.Marshal(pwds.Values())
			fmt.Println(string(bytes))
		} else {
			for _, k := range pwds.Keys() {
				fmt.Printf("%-30s %s\n", fmtx.Green.String(k), pwds.Get(k).ToString())
			}
		}
	}

	return nil
}
