package internal

import (
	"fmt"
	"testing"
)

func TestExternalTemplate(t *testing.T) {
	dir := "/Users/quanchao/code/go/src/jwzg/lanhu-demo/template"
	frame := "purger"
	files := GetExternalTemplateFiles(dir, frame)
	for _, file := range files {
		fmt.Println(file.Path)
	}
}
