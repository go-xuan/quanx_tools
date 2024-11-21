package copy

import (
	"github.com/go-xuan/quanx/net/ipx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Copy, "复制",
		flagx.BoolOption("blj_pwd", "蓝湖堡垒机SSH登录密码", false),
		flagx.BoolOption("ip", "本机IP", false),
	).SetExecutor(executor)
}

func executor() error {
	var text string
	if Command.GetOptionValue("lh_ssh_login").Bool() {
		text = "ssh quanchao@kicwhbttml.bastionhost.aliyuncs.com -p 60022"
	} else if Command.GetOptionValue("blj_pwd").Bool() {
		text = "371ADDd70c27_"
	} else if Command.GetOptionValue("ip").Bool() {
		text = ipx.GetLocalIP()
	} else {
		Command.OptionsHelp()
	}
	fmtx.Magenta.XPrintf("当前复制值：%s \n", text)
	if err := utils.CopyTobePasted(text); err != nil {
		return errorx.Wrap(err, "copy value to be pasted failed")
	}
	return nil
}
