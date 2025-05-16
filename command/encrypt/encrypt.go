package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/go-xuan/quanx/base/encodingx"
	"strings"
	"time"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/randx"

	"quanx_tools/command"
	"quanx_tools/common"
	"quanx_tools/common/enums"
	"quanx_tools/common/utils"
)

var (
	Command        = flagx.NewCommand(command.Encrypt, "根据公式加密")
	CryptoFunc     = enumx.NewStringEnum[string]()
	CryptoVariable = enumx.NewStringEnum[string]()
)

func init() {
	Command.AddOption(
		flagx.StringOption("formula", "加密公式", ""),
		flagx.StringOption("variables", "加密变量", ""),
		flagx.BoolOption("copy", "复制结果值", false),
	).SetExecutor(executor)

	CryptoFunc.
		Add("upper", `转大写，例如：-formula=upper(abc)，将字符串“abc”转为大写“ABC”`).
		Add("lower", `转小写，例如：-formula=upper(ABC)，将字符串“ABC”转为小写“abc”`).
		Add("reverse", `反转字符串，例如：-formula=reverse(abc)，将字符串“abc”转为“cba”`).
		Add("md5", `md5加密，例如：-formula=md5(abc)，将字符串“abc”进行md5加密`).
		Add("base64", `base64加密，例如：-formula=base64(abc)，将字符串“abc”进行base64加密`).
		Add("...", "扩展中，敬请期待...")

	CryptoVariable.
		Add("custom", `自定义常量值。例如：-variables="key=123"，将formula中的变量{key}赋值为自定义123`).
		Add("uuid", `特殊关键字。例如：-variables="key=uuid"，将formula中的变量{key}赋值为随机生成的UUID`).
		Add("timestamp", `特殊关键字。例如：-variables="key=timestamp"，将formula中的变量{key}赋值为当前秒级时间戳"`).
		Add("...", "扩展中，敬请期待...")
}

func executor() error {
	formula := Command.GetOptionValue("formula").String()
	if formula == "" {
		Command.OptionsHelp()
		fmt.Println("formula加密公式可根据不同需求进行自定义，完整加密公式由“加密函数+加密文本+加密变量”组成")
		fmt.Println(fmtx.Cyan.String("加密函数："), `加密函数在formula中使用为：“func1(...func2(...)...)”，加密函数支持嵌套。例如：-formula=md5(base64(abc))，先将文本“abc”进行base64加密，然后再进行md5加密`)
		fmt.Println(fmtx.Cyan.String("加密文本："), `加密文本可拼接在formula中的任意位置。例如：-formula=md5(abc_{app_id})，先将变量{app_id}拼接前缀“abc_”，然后再进行md5加密`)
		fmt.Println(fmtx.Cyan.String("加密变量："), `加密变量需要在formula中使用“{}”进行标识，并在variables中使用键值对进行赋值`)
		fmt.Println()

		fmt.Println(fmtx.Cyan.String("加密函数说明："))
		enums.Print(fmtx.Green, CryptoFunc)
		fmt.Println()

		fmt.Println(fmtx.Cyan.String("加密变量说明："))
		enums.Print(fmtx.Green, CryptoVariable)
		return nil
	}

	fmt.Println("加密公式: ", formula)
	variables := Command.GetOptionValue("variables").String()
	fmt.Println("加密变量: ", variables)

	var oldnew []string
	oldnew = append(oldnew, "{uuid}", randx.UUID())
	oldnew = append(oldnew, "{timestamp}", anyx.Int64Value(time.Now().Unix()).String())
	if params := stringx.ParseUrlParams(variables); len(params) > 0 {
		for k, v := range params {
			switch strings.ToLower(v) {
			case common.Uuid:
				v = randx.UUID()
			case common.Timestamp:
				v = anyx.Int64Value(time.Now().Unix()).String()
			}
			oldnew = append(oldnew, "{"+k+"}", v)
		}
	}
	formula = strings.NewReplacer(oldnew...).Replace(formula)
	fmt.Println("实际加密公式: ", formula)
	var result = doEncrypt(formula)
	fmt.Println("加密结果: ", fmtx.Magenta.String(result))
	// 开启复制
	if Command.GetOptionValue("copy").Bool() {
		if err := utils.WriteToClipboard(result); err != nil {
			return errorx.Wrap(err, "复制值到待粘贴失败")
		}
		fmtx.Magenta.Xprintf("当前值%s已复制到粘贴板\n", result)
	}
	return nil
}

// doEncrypt 执行加密
func doEncrypt(formula string) string {
	if funcName, start, end := getEncryptFuncAndIndex(formula); start > 0 {
		text := formula[start+1 : end]
		text = doEncrypt(text)
		switch funcName {
		case "upper":
			text = stringx.ToUpperCamel(text)
		case "lower":
			text = stringx.ToLowerCamel(text)
		case "reverse":
			text = stringx.Reverse(text)
		case "md5":
			text = encodingx.Hash(md5.New()).Encode([]byte(text))
		case "sha1":
			text = encodingx.Hash(sha1.New()).Encode([]byte(text))
		case "sha256":
			text = encodingx.Hash(sha256.New()).Encode([]byte(text))
		case "base64":
			text = encodingx.Base64(true).Encode([]byte(text))
		}
		start -= len(funcName)
		return formula[:start] + text + formula[end+1:]
	}
	return formula
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
