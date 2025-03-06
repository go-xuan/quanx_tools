package tts

import (
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/flagx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.TTS, "文本转语音（仅适用windows）",
		flagx.StringOption("content", "语音内容", "哈哈哈"),
	).SetExecutor(executor)
}

func executor() error {
	var content = Command.GetOptionValue("content").String()
	var output = filepath.Join(command.TTS, content+".mp3")
	if err := ole.CoInitialize(0); err != nil {
		return errorx.Wrap(err, "初始化COM库失败")
	}
	unknown, _ := oleutil.CreateObject("SAPI.SpVoice")
	voice, _ := unknown.QueryInterface(ole.IID_IDispatch)
	saveFile, _ := oleutil.CreateObject("SAPI.SpFileStream")
	ff, _ := saveFile.QueryInterface(ole.IID_IDispatch)
	// 打开wav文件
	_, _ = oleutil.CallMethod(ff, "Open", output, 3, true)
	// 设置voice的AudioOutputStream属性，必须是PutPropertyRef，如果是PutProperty就无法生效
	_, _ = oleutil.PutPropertyRef(voice, "AudioOutputStream", ff)
	// 设置语速
	_, _ = oleutil.PutProperty(voice, "Rate", 0)
	// 设置音量
	_, _ = oleutil.PutProperty(voice, "Volume", 400)
	// 说话
	_, _ = oleutil.CallMethod(voice, "Speak", content)
	// 停止说话
	//_, _ = oleutil.CallMethod(voice, "Pause")
	// 恢复说话
	//_, _ = oleutil.CallMethod(voice, "Resume")
	// 等待结束
	_, _ = oleutil.CallMethod(voice, "WaitUntilDone", 1000000)
	// 关闭文件
	_, _ = oleutil.CallMethod(ff, "Close")
	ff.Release()
	voice.Release()
	ole.CoUninitialize()
	return nil
}
