package copy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/base/ipx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/types/timex"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var (
	Command = flagx.NewCommand(command.Copy, "复制记录（可记录复制文本以备后续使用）")
	values  = enumx.NewStringEnum[*Value]()
	path    string
)

type Value struct {
	Key     string    `json:"key"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func (v *Value) Marshal() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v *Value) Unmarshal(j string) {
	_ = json.Unmarshal([]byte(j), v)
}

func (v *Value) Save(path string) error {
	return filex.WriteFileString(path, v.Marshal()+"\n", filex.Append)
}

func init() {
	Command.AddOption(
		flagx.BoolOption("list", "查看复制记录", false),
		flagx.BoolOption("clear", "清空复制记录", false),
	)
	Command.SetExecutor(executor)
	path = filepath.Join(os.Getenv("GOPATH"), "copy.yaml")
	if lines, _ := filex.ReadFileLine(path); len(lines) == 0 {
		v1 := &Value{Key: "1", Content: ipx.GetLocalIP(), Time: time.Now()}
		v2 := &Value{Key: "2", Content: "371ADDd70c27_", Time: time.Now()}
		v1.Save(path)
		v2.Save(path)
		values.Add(v1.Key, v1)
		values.Add(v2.Key, v2)
	} else {
		for _, line := range lines {
			value := &Value{}
			value.Unmarshal(line)
			values.Add(value.Key, value)
		}
	}
}

func Print(enum *enumx.Enum[string, *Value]) {
	for i := enum.Len(); i > 0; i-- {
		key := enum.Keys()[i-1]
		value := enum.Get(key)
		fmt.Printf("%s | %10s ==> %s\n",
			value.Time.Format(timex.TimeFmt),
			fmtx.Magenta.String(key),
			fmtx.Yellow.String(value.Content),
		)
	}
}

func executor() error {
	// 获取第一个参数
	arg := Command.GetArg(0)
	if Command.GetOptionValue("list").Bool() || arg == "list" {
		Print(values)
		return nil
	} else if Command.GetOptionValue("clear").Bool() || arg == "clear" {
		filex.Clear(path)
		if key := Command.GetArg(1); key != "" {
			values.Remove(key)
			for i, k := range values.Keys() {
				v := values.Get(k)
				v.Key = strconv.Itoa(i + 1)
				_ = v.Save(path)
			}
			fmtx.Red.Println("记录已删除")
		} else {
			values.Clear()
			fmtx.Red.Println("记录已全部清空")
		}
		return nil
	} else if arg == "" {
		// 获取当前粘贴板中复制内容
		if content, err := utils.ReadFromClipboard(); content != "" && err == nil {
			value := &Value{
				Key:     strconv.Itoa(values.Len() + 1),
				Content: content,
				Time:    time.Now(),
			}
			if err = value.Save(path); err != nil {
				return errorx.Wrap(err, "保存记录失败")
			} else {
				fmtx.Magenta.Xprintf("%s已保存记录\n", content)
			}
		}
	} else {
		if value := values.Get(arg); value != nil {
			text := values.Get(arg).Content
			if err := utils.WriteToClipboard(text); err != nil {
				return errorx.Wrap(err, "复制到粘贴板失败")
			}
			fmtx.Magenta.Xprintf("%s已复制到粘贴板\n", text)
		} else {
			fmtx.Magenta.Xprintf("%s暂无可复制值\n", arg)
		}
	}
	return nil
}
