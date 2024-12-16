package sql_fmt

import (
	"fmt"
	"path/filepath"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/sqlx/beautify"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.SqlFmt, "SQL格式化工具",
		flagx.StringOption("path", "SQL文件路径", ""),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetExecutor(executor)
}
func executor() error {
	if args := Command.GetArgs(); len(args) > 0 && args[0] == "-h" {
		Command.OptionsHelp()
	}
	var sql, outputPath, ifCopy = "", "beautify.sql", false
	if path := Command.GetOptionValue("path").String(); path != "" {
		path = filex.Pwd(path)
		fmt.Println("目标SQL文件：", path)
		if bytes, err := filex.ReadFile(path); err != nil {
			if !filex.Exists(path) {
				_ = filex.Create(path)
				fmt.Println("请在此SQL文件输入需要格式的SQL：", path)
				return err
			}
			return errorx.Wrap(err, "读取SQL文件失败")
		} else if len(bytes) == 0 {
			fmt.Println("请在此SQL文件输入需要格式的SQL：", path)
			return err
		} else {
			sql = string(bytes)
			var dir, name, suffix = filex.Analyse(path)
			outputPath = filepath.Join(dir, fmt.Sprintf("%s_fmt%s", name, suffix))
			ifCopy = Command.GetOptionValue("copy").Bool()
		}
	} else if content, err := utils.ReadFromClipboard(); content != "" && err == nil {
		sql, ifCopy = content, true
	} else {
		return errorx.Wrap(err, "获取SQL失败")
	}
	if len(sql) > 20 {
		// 美化sql
		var beautifySql = beautify.Parse(sql).Beautify()
		fmt.Println("格式化SQL:")
		fmtx.Green.Println(beautifySql)

		// 输出到粘贴板或者文件
		if ifCopy {
			if err := utils.WriteToClipboard(beautifySql); err != nil {
				return errorx.Wrap(err, "复制SQL到粘贴板失败")
			}
		} else {
			fmt.Println("写入SQL文件:", outputPath)
			if err := filex.WriteFileString(outputPath, beautifySql); err != nil {
				return errorx.Wrap(err, "写入SQL文件失败")
			}
		}
	}
	return nil
}
