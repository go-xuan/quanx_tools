package qr_code

import (
	"fmt"
	"github.com/go-xuan/quanx/utils/treex"
	"image/color"
	"path/filepath"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/utils/idx"
	"github.com/skip2/go-qrcode"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command = flagx.NewCommand(command.QrCode, "生成二维码")

func init() {
	Command.AddOption(
		flagx.StringOption("content", "二维码内容", ""),
		flagx.IntOption("size", "二维码大小", 900),
		flagx.BoolOption("copy", "复制结果值", false),
	).SetExecutor(executor)
}

func executor() error {
	content := Command.GetOptionValue("content").String()
	size := Command.GetOptionValue("size").Int()
	if content == "" {
		Command.OptionsHelp()
		return nil
	}
	content = treex.Trie().Desensitize(content)
	var name = idx.SnowFlake().String()
	path := filepath.Join(command.QrCode, name+".png")
	filex.CreateIfNotExist(path)
	if err := qrcode.WriteColorFile(
		content,        // 文本内容
		qrcode.Highest, // 级别
		size,           // 尺寸
		color.Black,    // 背景颜色
		color.White,    // 前景颜色
		path,           // 保存路径
	); err != nil {
		return errorx.Wrap(err, "生成二维码失败")
	}
	path = filex.Pwd(path)
	fmt.Println("二维码已保存至", fmtx.Yellow.String(path))
	if Command.GetOptionValue("copy").Bool() {
		if err := utils.WriteToClipboard(path); err != nil {
			return errorx.Wrap(err, "复制值二维码文件路径失败")
		}
		fmtx.Magenta.Xprintf("当前值%s已复制到粘贴板\n", path)
	}
	return nil
}
