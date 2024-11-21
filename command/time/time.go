package time

import (
	"fmt"
	"time"

	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/timex"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.TimeParse, "时间解析器",
		flagx.StringOption("parse", "待解析时间字符串，格式为YYYY-MM-DD hh:mm:ss", ""),
	).SetExecutor(executor)
}

// 时间解析器
func executor() error {
	parseTime := Command.GetOptionValue("parse").String()
	now := time.Now()
	if parseTime != "" {
		input := timex.ParseDateOrTime(parseTime)
		fmtx.Cyan.XPrintf("输入的时间是%s，时间戳是%v，%s年，星期%s，距离今天%v年/%v月/%v天\n",
			parseTime,
			input.UnixMilli(),
			timex.ShengXiao(input.Year()),
			timex.WeekdayCn(input),
			timex.TimeDiff(input, now, timex.Year),
			timex.TimeDiff(input, now, timex.Month),
			timex.TimeDiff(input, now, timex.Day),
		)
	} else {
		fmtx.Green.XPrintf("可使用%s参数，解析任意时间字符串。", "-parse")
		fmt.Printf(`例如：-parse="%s"`, timex.NowString())
		fmt.Println()
	}

	fmtx.Magenta.XPrintf("现在的时间是%s，时间戳是%v，%s年，星期%s，今天已过%v秒，今年已过%v天\n",
		now.Format(timex.TimeFmt),
		now.UnixMilli(),
		timex.ShengXiao(now.Year()),
		timex.WeekdayCn(now),
		timex.TimeDiff(timex.DateStart(now), now, timex.Second),
		now.YearDay())
	return nil
}
