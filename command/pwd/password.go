package pwd

import (
	"fmt"
	"io/fs"

	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
	"quanx_tools/command/pwd/conf"
	"quanx_tools/common/enums"
)

var (
	Command      = flagx.NewCommand(command.Password, "密码")
	PgsqlPwdEnum = enumx.NewStringEnum[string]()
	MysqlPwdEnum = enumx.NewStringEnum[string]()
	RedisPwdEnum = enumx.NewStringEnum[string]()
	MongoPwdEnum = enumx.NewStringEnum[string]()
)

type PWD struct {
	Type     string `yaml:"type" json:"type"`         // 数据库类型
	Env      string `yaml:"env" json:"env"`           // 环境
	Host     string `yaml:"host" json:"host"`         // host
	Port     string `yaml:"port" json:"port"`         // 端口
	Username string `yaml:"username" json:"username"` // 用户名
	Password string `yaml:"password" json:"password"` // 密码
}

func init() {
	Command.AddOption(
		flagx.BoolOption("pgsql", "pgsql密码", false),
		flagx.BoolOption("mysql", "mysql密码", false),
		flagx.BoolOption("redis", "redis密码", false),
		flagx.BoolOption("mongo", "mongo密码", false),
		flagx.StringOption("env", "环境", ""),
	).SetExecutor(executor)

	var enumMap = make(map[string]*enumx.StringEnum[string])
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
			var pwds []*PWD
			if err = marshalx.Apply(path).Unmarshal(content, &pwds); err != nil {
				return err
			}
			_, name, _ := filex.Analyse(path)
			enum := enumMap[name]
			for _, pwd := range pwds {
				enum.Add(pwd.Env, fmtx.Green.XSPrintf(
					"host: %-16s port: %s username: %-10s password: %s",
					pwd.Host, pwd.Port, pwd.Username, pwd.Password))
			}
		}
		return nil
	})
}

func executor() error {
	var enum *enumx.StringEnum[string]
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

	if enum != nil {
		env := Command.GetOptionValue("env").String()
		if pwd := enum.Get(env); pwd != "" {
			fmt.Println(pwd)
		} else {
			enums.Print(fmtx.Green, enum)
		}
	}
	return nil
}
