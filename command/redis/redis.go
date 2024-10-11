package redis

import (
	"context"
	"time"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/redisx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Redis, "redis删除工具",
		flagx.BoolOption("delete", "删除", false),
		flagx.BoolOption("set", "更新", false),
		flagx.StringOption("host", "主机", "127.0.0.1"),
		flagx.IntOption("port", "端口", 6379),
		flagx.StringOption("password", "密码", "admin@123"),
		flagx.IntOption("db", "数据库", 0),
		flagx.StringOption("key", "操作key", ""),
		flagx.StringOption("value", "操作值", ""),
	).SetHandler(handler)
}

func handler() error {
	deleteKey := Command.GetOptionValue("key").String()
	if deleteKey == "" {
		return nil
	}
	host := Command.GetOptionValue("host").String()
	port := Command.GetOptionValue("host").Int()
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
		redisx.Client().Del(ctx, deleteKey)
	} else if Command.GetOptionValue("set").Bool() {
		value := Command.GetOptionValue("value").String()
		redisx.Client().Set(ctx, deleteKey, value, time.Minute)
	}
	if err := redisx.Client().Close(); err != nil {
		return errorx.Wrap(err, "关闭redis连接失败")
	}
	return nil
}
