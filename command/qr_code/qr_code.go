package qr_code

import (
	"image/color"
	"path/filepath"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/utils/idx"
	"github.com/go-xuan/quanx/utils/treex"
	"github.com/skip2/go-qrcode"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.QrCode, "生成二维码",
		flagx.StringOption("content", "二维码内容", "123"),
		flagx.IntOption("size", "二维码大小", 600),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetHandler(handler)
}

func handler() error {
	content := Command.GetOptionValue("content").String()
	size := Command.GetOptionValue("size").Int(600)
	if content != "" {
		content = treex.Trie().Desensitize(content)
	}
	var name = idx.SnowFlake().String()
	path := filepath.Join("qrCode", name+".png")
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
	fmtx.Blue.XPrintf("二维码保存至：%s", path)
	if Command.GetOptionValue("copy").Bool() {
		if err := utils.CopyTobePasted(path); err != nil {
			return errorx.Wrap(err, "复制值二维码文件路径失败")
		}
	}
	return nil
}
