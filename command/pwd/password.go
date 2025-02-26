package pwd

import (
	"fmt"
	"io/fs"
	"strconv"

	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
	"quanx_tools/command/pwd/conf"
)

var (
	Command      = flagx.NewCommand(command.Password, "密码")
	PgsqlPwdEnum = enumx.NewStringEnum[*Password]()
	MysqlPwdEnum = enumx.NewStringEnum[*Password]()
	RedisPwdEnum = enumx.NewStringEnum[*Password]()
	MongoPwdEnum = enumx.NewStringEnum[*Password]()
)

type Password struct {
	Type     string `yaml:"type" json:"type"`         // 数据库类型
	Env      string `yaml:"env" json:"env"`           // 环境
	Host     string `yaml:"host" json:"host"`         // host
	Port     int    `yaml:"port" json:"port"`         // 端口
	Database string `yaml:"database" json:"database"` // 数据库
	Username string `yaml:"username" json:"username"` // 用户名
	Password string `yaml:"password" json:"password"` // 密码
}

func (p *Password) ToString() string {
	return fmtx.Yellow.XSPrintf(
		"host: %s port: %s database: %s username: %s password: %s",
		stringx.Fill(p.Host, " ", 15),
		strconv.Itoa(p.Port),
		stringx.Fill(p.Database, " ", 15),
		stringx.Fill(p.Username, " ", 10),
		p.Password)
}

func init() {
	Command.AddOption(
		flagx.BoolOption("pgsql", "获取pgsql连接信息", false),
		flagx.BoolOption("mysql", "获取mysql连接信息", false),
		flagx.BoolOption("redis", "获取redis连接信息", false),
		flagx.BoolOption("mongo", "获取mongo连接信息", false),
		flagx.BoolOption("json", "json格式化", false),
		flagx.BoolOption("yaml", "yaml格式化", false),
		flagx.StringOption("env", "选择环境", ""),
	).SetExecutor(executor)

	var enumMap = make(map[string]*enumx.StringEnum[*Password])
	enumMap["pgsql"] = PgsqlPwdEnum
	enumMap["mysql"] = MysqlPwdEnum
	enumMap["redis"] = RedisPwdEnum
	enumMap["mongo"] = MongoPwdEnum

	_ = fs.WalkDir(conf.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filex.GetSuffix(path) != "go" {
			var content []byte
			if content, err = conf.FS.ReadFile(path); err != nil {
				return err
			}
			var pwds []*Password
			if err = marshalx.Apply(path).Unmarshal(content, &pwds); err != nil {
				return err
			}
			_, name, _ := filex.Analyse(path)
			enum := enumMap[name]
			for _, pwd := range pwds {
				enum.Add(pwd.Env, pwd)
			}
		}
		return nil
	})
}

func executor() error {
	var enum *enumx.StringEnum[*Password]
	if Command.GetOptionValue("pgsql").Bool() {
		enum = PgsqlPwdEnum
	} else if Command.GetOptionValue("mysql").Bool() {
		enum = MysqlPwdEnum
	} else if Command.GetOptionValue("redis").Bool() {
		enum = RedisPwdEnum
	} else if Command.GetOptionValue("mongo").Bool() {
		enum = MongoPwdEnum
	} else {
		Command.OptionsHelp()
		return nil
	}

	// 序列化方式
	var method marshalx.Method
	if Command.GetOptionValue("json").Bool() {
		method = marshalx.Json{Indent: "    "}
	} else if Command.GetOptionValue("yaml").Bool() {
		method = marshalx.Yaml{}
	}

	if enum != nil {
		if pwd := enum.Get(Command.GetOptionValue("env").String()); pwd != nil {
			if method != nil {
				bytes, _ := method.Marshal(pwd)
				fmt.Println(string(bytes))
			} else {
				fmt.Println(pwd.ToString())
			}
		} else {
			if method != nil {
				bytes, _ := method.Marshal(enum.Values())
				fmt.Println(string(bytes))
			} else {
				for _, k := range enum.Keys() {
					fmt.Printf("%-30s %s\n", fmtx.Green.String(k), enum.Get(k).ToString())
				}
			}
		}
	}
	return nil
}
