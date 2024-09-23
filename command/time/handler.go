package time

import (
	"time"

	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/types/timex"
	"github.com/go-xuan/quanx/utils/fmtx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.TimeParse, "时间解析器",
		flagx.StringOption("parse", "待解析时间字符串，格式为YYYY-MM-DD hh:mm:ss", "")).
		SetHandler(handler)
}

// 时间解析器
func handler() error {
	parseTime := Command.GetOptionValue("parse").String()
	now := time.Now()
	if parseTime != "" {
		input := timex.Parse(parseTime)
		fmtx.Cyan.XPrintf("输入的时间是%s，时间戳是%v，%s年，星期%s，距离今天%v天、%v月、%v年",
			parseTime,
			input.UnixMilli(),
			timex.ShengXiao(input.Year()),
			timex.WeekdayCn(input),
			timex.TimeDiff(input, now, timex.Day),
			timex.TimeDiff(input, now, timex.Month),
			timex.TimeDiff(input, now, timex.Year))
	}
	fmtx.Magenta.XPrintf("现在的时间是%s，时间戳是%v，%s年，星期%s，今天已过%v秒，今年已过%v天",
		now.Format(timex.TimeFmt),
		now.UnixMilli(),
		timex.ShengXiao(now.Year()),
		timex.WeekdayCn(now),
		timex.TimeDiff(timex.DateStart(now), now, timex.Second),
		now.YearDay())
	return nil
}
