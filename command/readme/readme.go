package readme

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/enumx"

	"quanx_tools/command"
)

//go:embed *
var FS embed.FS

var (
	Command = flagx.NewCommand(command.Readme, "README")
	allData = enumx.NewStringEnum[*enumx.Enum[string, []byte]]()
)

func init() {
	_ = fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filex.GetSuffix(path) != "go" {
			group, name, _ := filex.Analyse(path)
			group = strings.Trim(group, "/")
			var data *enumx.Enum[string, []byte]
			if data = allData.Get(group); data == nil {
				data = enumx.NewStringEnum[[]byte]()
			}
			var content []byte
			if content, err = FS.ReadFile(path); err != nil {
				return err
			}
			data.Add(name, content)
			allData.Add(group, data)
		}
		return nil
	})

	// 给readme命令添加可选项
	if groups := allData.Keys(); len(groups) > 0 {
		for _, group := range groups {
			Command.AddOption(flagx.StringOption(group, "可选分组", ""))
		}
	}

	Command.AddOption(flagx.BoolOption("copy", "复制结果值", false))
	Command.SetExecutor(executor)
}

func executor() error {
	var group string
	var groupData *enumx.Enum[string, []byte]
	if group = Command.GetArg(0); group != "" {
		group = strings.TrimPrefix(group, "-")
		groupData = allData.Get(group)
	}
	if groupData == nil {
		fmtx.Red.Println("请指定一个可选分组：")
		for _, key := range allData.Keys() {
			fmtx.Magenta.Println(key)
		}
		return nil
	} else if groupData.Len() == 1 {
		fmtx.Green.Println(string(groupData.Values()[0]))
		return nil
	} else if name := Command.GetArg(1); name != "" {
		if content := groupData.Get(name); content != nil {
			fmtx.Green.Println(string(content))
			return nil
		}
	} else {
		fmtx.Red.Xprintf("请指定%s分组下其中一项：\n", group)
		for _, key := range groupData.Keys() {
			fmtx.Magenta.Println(key)
		}
	}
	return nil
}
