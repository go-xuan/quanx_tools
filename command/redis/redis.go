package redis

import (
	"context"
	"time"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/redisx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"

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
	key := Command.GetOptionValue("key").String()
	if key == "" {
		fmtx.Red.Println("key is empty")
		return nil
	}
	host := Command.GetOptionValue("host").String()
	port := Command.GetOptionValue("port").Int()
	password := Command.GetOptionValue("password").String()
	db := Command.GetOptionValue("db").Int()
	// 初始化
	if err := configx.Execute(&redisx.Config{
		Source:   "default",
		Enable:   true,
		Host:     host,
		Port:     port,
		Password: password,
		Database: db,
		PoolSize: 15,
	}); err != nil {
		return errorx.Wrap(err, "初始化redis连接失败")
	}

	var ctx = context.TODO()
	if Command.GetOptionValue("delete").Bool() {
		redisx.GetClient().Del(ctx, key)
	} else if Command.GetOptionValue("set").Bool() {
		value := Command.GetOptionValue("value").String()
		redisx.GetClient().Set(ctx, key, value, time.Minute)
	} else if Command.GetOptionValue("get").Bool() {
		value := redisx.GetClient().Get(ctx, key)
		fmtx.Cyan.Printf("the value of %s is: %s", key, value)
	}
	if err := redisx.GetClient().Close(); err != nil {
		return errorx.Wrap(err, "关闭redis连接失败")
	}
	return nil
}
