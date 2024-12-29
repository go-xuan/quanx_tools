package utils

import (
	"bytes"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/execx"
)

func WriteToClipboard(text string) error {
	var command string
	switch runtime.GOOS {
	case "windows":
		command = "clip"
	case "darwin":
		command = "pbcopy"
	default:
		return nil
	}
	if _, _, err := execx.Command(command).Stdin(bytes.NewBufferString(text)).Run(); err != nil {
		return errorx.Wrap(err, "copy value to clipboard failed")
	}
	return nil
}

func ReadFromClipboard() (content string, err error) {
	switch runtime.GOOS {
	case "windows":
		content, err = clipboard.ReadAll()
	case "darwin":
		content, _, err = execx.Command("pbpaste").Run()
	default:
		err = errorx.New("unknown os")
	}
	if err != nil {
		err = errorx.Wrap(err, "failed to read from clipboard")
	}
	return
}
