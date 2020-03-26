package app

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/king19800105/go-kit-ddd/pkg/config"
	"github.com/king19800105/go-kit-ddd/pkg/storage"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 服务启动逻辑
func Run() error {
	// 初始化指令
	pflag.StringP("http", "h", ":8080", "http listen address")
	pflag.StringP("env", "e", "dev", "env set")

	// 加载配置
	cfg, err := config.Viperize()
	if nil != err {
		return err
	}

	// 日志设置
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	// 数据库设置
	lifeTime := time.Duration(cfg.GetInt64("database.max_life_time")) * time.Second
	storage.WithLog(cfg.GetBool("database.use_log"))
	storage.WithPrefix(cfg.GetString("database.prefix"))
	storage.WithMaxLifeTime(lifeTime)
	storage.WithMaxIdleConns(cfg.GetInt("max_idle_conns"))
	storage.WithMaxOpenConns(cfg.GetInt("max_open_conns"))
	db, err := storage.NewDB(
		cfg.GetString("database.dialect"),
		cfg.GetString("database.host"),
		cfg.GetString("database.username"),
		cfg.GetString("database.password"),
		cfg.GetString("database.name"))
	if nil != err {
		return err
	}

	// 全局ctx对象
	_ = context.Background()
	// 错误管道创建
	errCh := make(chan error)

	// 启动http服务监听
	go func() {
		defer func() {
			level.Info(logger).Log("msg", "base operating http service ended")
			db.Close()
		}()
		level.Info(logger).Log("msg", "base operating http service started")
	}()

	// 信号量监听
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("信号量中断，%s", <-c)
	}()

	return <-errCh
}
