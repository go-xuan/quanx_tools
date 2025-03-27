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
	Command = flagx.NewCommand(command.Copy, "复制板")
	Data    = enumx.NewStringEnum[*Value]()
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
	Command.SetExecutor(executor)
	path = filepath.Join(os.Getenv("GOPATH"), "copy.yaml")
	if lines, _ := filex.ReadFileLine(path); len(lines) == 0 {
		v1 := &Value{Key: "bljmm", Content: "371ADDd70c27_", Time: time.Now()}
		v2 := &Value{Key: "ip", Content: ipx.GetLocalIP(), Time: time.Now()}
		v1.Save(path)
		v2.Save(path)
		Data.Add(v1.Key, v1)
		Data.Add(v2.Key, v2)
	} else {
		for _, line := range lines {
			value := &Value{}
			value.Unmarshal(line)
			Data.Add(value.Key, value)
		}
	}
}

func Print(color fmtx.Color, enum *enumx.Enum[string, *Value]) {
	for i := enum.Len(); i > 0; i-- {
		key := enum.Keys()[i-1]
		value := enum.Get(key)
		fmt.Printf("%-30s %20s %s\n", color.String(key), value.Time.Format(timex.TimeFmt), value.Content)
	}
}

func executor() error {
	arg := Command.GetArg(0)
	switch arg {
	case "list":
		Print(fmtx.Magenta, Data)
	case "clear":
		Data.Clear()
		filex.Clear(path)
	case "":
		if content, err := utils.ReadFromClipboard(); content != "" && err == nil {
			v := &Value{
				Key:     strconv.Itoa(len(Data.Keys()) + 1),
				Content: content,
				Time:    time.Now(),
			}
			if err = v.Save(path); err != nil {
				return errorx.Wrap(err, "保存粘贴板失败")
			} else {
				fmtx.Magenta.XPrintf("%s已保存至粘贴板\n", content)
			}
		}
	default:
		if value := Data.Get(arg); value != nil {
			text := Data.Get(arg).Content
			if err := utils.WriteToClipboard(text); err != nil {
				return errorx.Wrap(err, "复制到粘贴板失败")
			}
			fmtx.Magenta.XPrintf("%s已复制到粘贴板\n", text)
		} else {
			fmtx.Magenta.XPrintf("%s暂无可复制值\n", arg)
		}
	}
	return nil
}
