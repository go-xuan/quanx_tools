package tts

import (
	"os/exec"
	
	"github.com/go-xuan/quanx/base/errorx"
)

func ttsDarwin(content string) error {
	if err := exec.Command("say", "-o", content+".aiff", content).Run(); err != nil {
		return errorx.Wrap(err, "文本转语音失败")
	}
	return nil
}
