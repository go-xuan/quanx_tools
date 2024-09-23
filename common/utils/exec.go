package utils

import (
	"bytes"
	"runtime"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/execx"
)

func CopyTobePasted(text string) error {
	var command string
	switch runtime.GOOS {
	case "windows":
		command = "clip"
	case "darwin":
		command = "pbcopy"
	default:
		return nil
	}
	if _, err := execx.ExecCommand(command, bytes.NewBufferString(text)); err != nil {
		return errorx.Wrap(err, "copy value to clipboard failed")
	}
	return nil
}
