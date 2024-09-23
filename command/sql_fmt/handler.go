package sql_fmt

import (
	"fmt"
	"path/filepath"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/utils/fmtx"
	"github.com/go-xuan/quanx/utils/sqlx"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.SqlFmt, "SQL格式化工具",
		flagx.StringOption("path", "SQL文件路径", "fmt.sql"),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetHandler(handler)
}

func handler() error {
	inputPath := Command.GetOptionValue("path").String()
	fmt.Println("输入SQL格式化文件：", inputPath)
	if bytes, err := filex.ReadFile(inputPath); err != nil {
		fmtSql := sqlx.Format(string(bytes)).String()
		fmt.Println("格式化SQL:")
		fmtx.Green.Println(fmtSql)
		if Command.GetOptionValue("copy").Bool() {
			if err = utils.CopyTobePasted(fmtSql); err != nil {
				return errorx.Wrap(err, "复制值到待粘贴内容失败")
			}
		} else {
			dir, name, suffix := filex.Analyse(inputPath)
			outputPath := filepath.Join(dir, fmt.Sprintf("%s_fmt%s", name, suffix))
			fmt.Println("输出SQL格式化文件:", outputPath)
			if err = filex.WriteFile(outputPath, fmtSql); err != nil {
				return errorx.Wrap(err, "写入sql文件失败")
			}
		}
	} else {
		return errorx.Wrap(err, "读取sql文件失败")
	}
	return nil
}
