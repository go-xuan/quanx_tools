package redis

import (
	"context"
	"time"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/extra/redisx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Redis, "redis工具",
		flagx.BoolOption("delete", "删除", false),
		flagx.BoolOption("set", "更新", false),
		flagx.BoolOption("get", "查询", false),
		flagx.StringOption("host", "主机", "127.0.0.1"),
		flagx.IntOption("port", "端口", 6379),
		flagx.StringOption("password", "密码", "admin@123"),
		flagx.IntOption("db", "数据库", 0),
		flagx.StringOption("key", "操作key", ""),
		flagx.StringOption("value", "操作值", ""),
	).SetExecutor(executor)
}

func executor() error {
	redis := &redisx.Config{
		Source:   "default",
		Enable:   true,
		Host:     Command.GetOptionValue("host").String(),
		Port:     Command.GetOptionValue("port").Int(),
		Password: Command.GetOptionValue("password").String(),
		Database: Command.GetOptionValue("db").Int(),
		PoolSize: 15,
	}
	// 初始化
	if err := redis.Execute(); err != nil {
		return errorx.Wrap(err, "初始化redis连接失败")
	}

	key := Command.GetOptionValue("key").String()
	if key == "" {
		fmtx.Red.Println("key is empty")
		return nil
	}

	var ctx = context.TODO()
	if Command.GetOptionValue("delete").Bool() {
		redisx.GetInstance().Del(ctx, key)
	} else if Command.GetOptionValue("set").Bool() {
		value := Command.GetOptionValue("value").String()
		redisx.GetInstance().Set(ctx, key, value, time.Minute)
	} else if Command.GetOptionValue("get").Bool() {
		value := redisx.GetInstance().Get(ctx, key)
		fmtx.Cyan.Printf("the value of %s is: %s", key, value)
	}
	if err := redisx.GetInstance().Close(); err != nil {
		return errorx.Wrap(err, "关闭redis连接失败")
	}
	return nil
}
