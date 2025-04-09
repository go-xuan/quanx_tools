package time

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/types/timex"

	"quanx_tools/command"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.TimeParse, "时间解析器",
		flagx.StringOption("str", "时间字符串，格式为YYYY-MM-DD hh:mm:ss", ""),
		flagx.Int64Option("unix", "时间戳", 0),
	).SetExecutor(executor)
}

// 时间解析器
func executor() error {
	var now = time.Now()
	fmtx.Magenta.Xprintf("现在的时间是%s，时间戳是%v，%s年，星期%s，今天已过%v秒，今年已过%v天\n",
		now.Format(timex.TimeFmt),
		now.UnixMilli(),
		timex.ShengXiao(now.Year()),
		timex.WeekdayCn(now),
		timex.TimeDiff(timex.DateStart(now), now, timex.Second),
		now.YearDay())

	var inputTime time.Time
	if str := Command.GetOptionValue("str").String(); str != "" {
		inputTime = timex.ParseDateOrTime(str)
	} else if unix := Command.GetOptionValue("unix").Int64(); unix > 0 {
		if unix > 1e12 {
			unix = unix / 1000
		}
		inputTime = time.Unix(unix, 0)
	} else if content, err := utils.ReadFromClipboard(); content != "" && err == nil {
		content = strings.TrimSpace(content)
		if inputTime = timex.ParseDateOrTime(content); inputTime.IsZero() {
			if unix = stringx.ParseInt64(content); unix > 0 {
				if unix > 1e12 {
					inputTime = time.UnixMilli(unix)
				} else {
					inputTime = time.Unix(unix, 0)
				}
			}
		}
	}
	if !inputTime.IsZero() {
		fmtStr, unixMilli := inputTime.Format("2006-01-02 15:04:05.999"), inputTime.UnixMilli()
		fmtx.Cyan.Xprintf("输入的时间是%s，时间戳是%v，%s年，星期%s，距离今天%v年，%v月，%v天\n",
			fmtStr,
			unixMilli,
			timex.ShengXiao(inputTime.Year()),
			timex.WeekdayCn(inputTime),
			timex.TimeDiff(inputTime, now, timex.Year),
			timex.TimeDiff(inputTime, now, timex.Month),
			timex.TimeDiff(inputTime, now, timex.Day),
		)
	} else {
		fmtx.Green.Xprintf("可使用%s参数，解析任意时间字符串。", "-str")
		fmt.Printf(`例如：-str="%s"`, timex.NowString())
		fmt.Println()
	}

	return nil
}
