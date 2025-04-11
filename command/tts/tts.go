package tts

import (
	"runtime"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/flagx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.TTS, "文本转语音（仅适用windows）").SetExecutor(executor)
}

func executor() error {
	if content := Command.GetArg(0); content != "" {
		switch runtime.GOOS {
		case `windows`:
			return ttsWindows(content)
		case `darwin`:
			return ttsDarwin(content)
		default:
			return errorx.New("unsupported platform")
		}
	}
	return errorx.New("content is empty")
}
