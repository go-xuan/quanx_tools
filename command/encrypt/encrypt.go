package encrypt

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/encryptx"
	"github.com/go-xuan/quanx/utils/randx"

	"quanx_tools/command"
	"quanx_tools/common"
	"quanx_tools/common/enums"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Encrypt, "根据公式加密",
		flagx.StringOption("formula", "加密公式", ""),
		flagx.StringOption("variables", "加密变量", ""),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetHandler(handler)
}

func handler() error {
	formula := Command.GetOptionValue("formula").String()

	fmt.Println("input formula: ", formula)
	if formula == "" {
		fmtx.Red.Println("formula加密公式可根据需求进行自定义，公式需由“加密函数”和“加密内容”两部分组成")
		fmt.Println(fmtx.Cyan.String("加密函数："), `格式为：“加密函数(加密内容)”，加密函数可嵌套。例如：formula=md5(base64(abc))， 将文本“abc”先进行base64加密之后再进行md5加密`)
		fmt.Println(fmtx.Cyan.String("加密内容："), `在formula中，加密内容可将“文本值”和“加密变量”拼接组合。例如：formula=md5(abc_{app_id})，变量{app_id}拼接“abc_”前缀之后再进行md5加密`)
		fmt.Println(fmtx.Cyan.String("加密变量："), `在formula中，引用变量需要使用“{}”进行标识，并需要在variables中对相应变量进行赋值`)
		fmt.Println()

		fmtx.Magenta.XPrintf(`%s函数使用：`, "formula")
		enums.Print(fmtx.Green, enums.CryptoFunc)
		fmt.Println()

		fmtx.Magenta.XPrintf(`%s变量使用：`, "variables")
		enums.Print(fmtx.Green, enums.CryptoVariable)
	}

	variables := Command.GetOptionValue("variables").String()
	fmt.Println("input variables: ", variables)
	params := stringx.ParseUrlParams(variables)
	if len(params) > 0 {
		var oldnew []string
		for k, v := range params {
			switch strings.ToLower(v) {
			case common.Uuid:
				v = randx.UUID()
			case common.Timestamp:
				v = anyx.Int64Value(time.Now().Unix()).String()
			}
			params[k] = v
			oldnew = append(oldnew, "{"+k+"}", v)
		}
		formula = strings.NewReplacer(oldnew...).Replace(formula)
	}
	fmt.Println("actual formula: ", formula)
	var result = doCrypto(formula)
	fmtx.Magenta.XPrintf("encrypt result: %s", result)
	// 开启复制
	if Command.GetOptionValue("copy").Bool() {
		if err := utils.CopyTobePasted(result); err != nil {
			return errorx.Wrap(err, "复制值到待粘贴失败")
		}
	}
	return nil
}

// doCrypto 执行加密
func doCrypto(formula string) (result string) {
	if funcName, start, end := getEncryptFuncAndIndex(formula); start > 0 {
		text := formula[start:end]
		text = doCrypto(text)
		switch funcName {
		case "upper":
			text = stringx.ToUpperCamel(text)
		case "lower":
			text = stringx.ToLowerCamel(text)
		case "reverse":
			text = stringx.Reverse(text)
		case "md5":
			text = encryptx.MD5(text)
		case "base64":
			text = encryptx.Base64Encode([]byte(text), true)
		}
		start -= len(funcName) + 1
		result = formula[:start] + text + formula[end+1:]
	} else {
		result = formula
	}
	return
}

// getEncryptFuncAndIndex 获取加密方法以及下标
func getEncryptFuncAndIndex(formula string) (string, int, int) {
	if start, end := stringx.Between(formula, "(", ")"); start < 0 || end < 0 {
		return "", 0, 0
	} else {
		funcName, _ := stringx.Contains(formula[:start], "md5", "upper", "lower", "base64", "reverse")
		return funcName, start, end
	}
}
