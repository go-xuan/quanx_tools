package pwd

import (
	"fmt"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/enumx"

	"quanx_tools/command"
	"quanx_tools/common/enums"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Password, "密码",
		flagx.BoolOption("pgsql", "pgsql密码", false),
		flagx.BoolOption("mysql", "mysql密码", false),
		flagx.BoolOption("redis", "redis密码", false),
		flagx.BoolOption("mongo", "mongo密码", false),
		flagx.StringOption("env", "环境", ""),
	).SetExecutor(executor)
}

func executor() error {
	env := Command.GetOptionValue("env").String()
	var enum *enumx.StringEnum[string]
	if Command.GetOptionValue("pgsql").Bool() {
		enum = enums.PgsqlPwdEnum
	} else if Command.GetOptionValue("mysql").Bool() {
		enum = enums.MysqlPwdEnum
	} else if Command.GetOptionValue("redis").Bool() {
		enum = enums.RedisPwdEnum
	} else if Command.GetOptionValue("mongo").Bool() {
		enum = enums.MongoPwdEnum
	} else {
		Command.OptionsHelp()
		return nil
	}
	if enum != nil {
		if pwd := enum.Get(env); pwd != "" {
			fmt.Println(pwd)
		} else {
			enums.Print(fmtx.Green, enum)
		}
	}
	return nil
}
