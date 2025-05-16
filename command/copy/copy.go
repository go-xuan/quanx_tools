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
	"github.com/go-xuan/quanx/base/osx"
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
	NO      string    `json:"no"`
	Tag     string    `json:"tag"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func (v *Value) Unmarshal(j string) {
	_ = json.Unmarshal([]byte(j), v)
}

func (v *Value) Save(path string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return errorx.Wrap(err, "value marshal error")
	}
	data = append(data, []byte("\n")...)
	return filex.WriteFile(path, data, filex.Append)
}

func init() {
	Command.AddOption(
		flagx.BoolOption("list", "查看复制记录", false),
		flagx.BoolOption("clear", "清理复制记录", false),
		flagx.BoolOption("save", "保存复制记录 | save [tag]", false),
		flagx.BoolOption("get", "获取复制记录 | get [no]", false),
		flagx.BoolOption("tag", "更新记录标签 | tag [no] [tag]", false),
	)
	Command.SetExecutor(executor)
	path = filepath.Join(os.Getenv("GOPATH"), "copy.yaml")
	if lines, _ := filex.ReadFileLine(path); len(lines) == 0 {
		v1 := &Value{NO: "1", Tag: "localhost", Content: osx.GetLocalIP(), Time: time.Now()}
		v2 := &Value{NO: "2", Tag: "baoleijimima", Content: "371ADDd70c27_", Time: time.Now()}
		v1.Save(path)
		v2.Save(path)
		values.Add(v1.NO, v1)
		values.Add(v2.NO, v2)
	} else {
		for _, line := range lines {
			value := &Value{}
			value.Unmarshal(line)
			values.Add(value.NO, value)
		}
	}
}

func executor() error {
	// 获取第一个参数
	arg0 := Command.GetArg(0)
	switch arg0 {
	case "list", "-list":
		return listValues()
	case "clear", "-clear":
		return clearValues()
	case "get", "-get":
		arg1 := Command.GetArg(1) // no
		return geValue(arg1)
	case "save", "-save":
		arg1 := Command.GetArg(1) // tag
		return saveValue(arg1)
	case "tag", "-tag":
		arg1 := Command.GetArg(1) // no
		arg2 := Command.GetArg(2) // tag
		return updateValueTag(arg1, arg2)
	default:
		return saveValue("")
	}
}

func listValues() error {
	fmtx.Magenta.Println(fmt.Sprintf("%-19s | %-5s | %-20s | %s", "CREATE_TIME", "NO", "TAG", "CONTENT"))
	for i := values.Len(); i > 0; i-- {
		key := values.Keys()[i-1]
		value := values.Get(key)
		fmt.Printf("%s | %-5s | %-20s | %s\n",
			value.Time.Format(timex.TimeFmt),
			value.NO,
			value.Tag,
			fmtx.Yellow.String(value.Content),
		)
	}
	return nil
}

// 清理复制记录
func clearValues() error {
	filex.Clear(path)
	if arg1 := Command.GetArg(1); arg1 != "" {
		values.Remove(arg1)
		for i, k := range values.Keys() {
			v := values.Get(k)
			v.NO = strconv.Itoa(i + 1)
			_ = v.Save(path)
		}
		fmtx.Red.Println("复制记录已删除")
	} else {
		values.Clear()
		fmtx.Red.Println("复制记录已清空")
	}
	return nil
}

// 获取复制记录
func geValue(no string) error {
	value := values.Get(no)
	if value == nil {
		fmtx.Magenta.Xprintf("%s暂无可复制值\n", no)
		return errorx.New("未找到对应记录")
	}
	text := values.Get(no).Content
	if err := utils.WriteToClipboard(text); err != nil {
		return errorx.Wrap(err, "复制到粘贴板失败")
	}
	fmtx.Magenta.Xprintf("%s已复制到粘贴板\n", text)
	return nil
}

// 获取当前粘贴板中复制内容并保存
func saveValue(tag string) error {
	content, err := utils.ReadFromClipboard()
	if err != nil || content == "" {
		return errorx.Wrap(err, "未能获取复制值")
	}
	value := &Value{
		NO:      strconv.Itoa(values.Len() + 1),
		Tag:     tag,
		Content: content,
		Time:    time.Now(),
	}
	if err = value.Save(path); err != nil {
		return errorx.Wrap(err, "保存记录失败")
	}
	fmtx.Magenta.Xprintf("%s已保存记录\n", content)
	return nil
}

// 获取当前粘贴板中复制内容并保存
func updateValueTag(no, tag string) error {
	filex.Clear(path)
	value := values.Get(no)
	if value == nil {
		return errorx.New("未找到相关记录")
	}
	oldTag := value.Tag
	value.Tag = tag
	for _, k := range values.Keys() {
		if err := values.Get(k).Save(path); err != nil {
			return errorx.Wrap(err, "保存记录失败")
		}
	}
	fmtx.Magenta.Xprintf("已更新%s的标签：%s ==> %s \n", no, oldTag, tag)
	return nil
}
