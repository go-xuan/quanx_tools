package sql_fmt

import (
	"fmt"
	"path/filepath"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/utils/sqlx"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.SqlFmt, "SQL格式化工具",
		flagx.StringOption("path", "SQL文件路径", ""),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetHandler(handler)
}

func handler() error {
	path := Command.GetOptionValue("path").String()
	if path == "" {
		fmtx.Green.XPrintf("可使用%s参数，指向需要格式化的sql文件\n", "-path")
		return nil
	}
	var inputPath = filex.Pwd(path)
	fmt.Println("目标SQL文件：", inputPath)
	if bytes, err := filex.ReadFile(inputPath); err != nil {
		if !filex.Exists(inputPath) {
			_ = filex.Create(inputPath)
			fmt.Println("请在此SQL文件输入需要格式的SQL：", inputPath)
			return nil
		}
		return errorx.Wrap(err, "读取SQL文件失败")
	} else if len(bytes) == 0 {
		fmt.Println("请在此SQL文件输入需要格式的SQL：", inputPath)
		return nil
	} else {
		var fmtSql = sqlx.Format(string(bytes)).String()
		fmt.Println("格式化SQL:")
		fmtx.Green.Println(fmtSql)
		if Command.GetOptionValue("copy").Bool() {
			if err = utils.CopyTobePasted(fmtSql); err != nil {
				return errorx.Wrap(err, "复制值到待粘贴内容失败")
			}
		} else {
			var dir, name, suffix = filex.Analyse(inputPath)
			var outputPath = filepath.Join(dir, fmt.Sprintf("%s_fmt%s", name, suffix))
			fmt.Println("写入SQL文件:", outputPath)
			if err = filex.WriteFile(outputPath, fmtSql); err != nil {
				return errorx.Wrap(err, "写入SQL文件失败")
			}
		}
	}
	return nil
}
