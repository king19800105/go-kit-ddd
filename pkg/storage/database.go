package storage

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type Option interface {
	apply(option *options)
}

type options struct {
	useLog  bool
	timeout string
	charset string
	prefix  string
	maxLife time.Duration
	maxIdle int
	maxOpen int
}

type optionFunc func(*options)

const dbConnect = "%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%s"

func (f optionFunc) apply(o *options) {
	f(o)
}

// 实例化 charset, prefix db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/user?charset=utf8&parseTime=True&loc=Local")
func NewDB(dialect, host, username, password, dbname string, opts ...Option) (db *gorm.DB, err error) {
	options := new(options)
	for _, o := range opts {
		o.apply(options)
	}

	if "" == options.charset {
		options.charset = "utf8"
	}

	if "" == options.timeout {
		options.timeout = "10ms"
	}

	connect := fmt.Sprintf(dbConnect, username, password, host, dbname, options.charset, options.timeout)
	db, err = gorm.Open(dialect, connect)
	if nil != err {
		return
	}

	initOptions(db, options)
	return
}

// 日志使用设置
func WithLog(ok bool) Option {
	return optionFunc(func(o *options) {
		o.useLog = ok
	})
}

// 超时设置
func WithTimeout(t string) Option {
	return optionFunc(func(o *options) {
		o.timeout = t
	})
}

// 字符集设置
func WithCharset(c string) Option {
	return optionFunc(func(o *options) {
		o.charset = c
	})
}

// 表前缀设置
func WithPrefix(pre string) Option {
	return optionFunc(func(o *options) {
		o.prefix = pre
	})
}

// 链接复用周期
func WithMaxLifeTime(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.maxLife = t
	})
}

// 空间连接数
func WithMaxIdleConns(cnt int) Option {
	return optionFunc(func(o *options) {
		o.maxIdle = cnt
	})
}

// 满载连接数
func WithMaxOpenConns(cnt int) Option {
	return optionFunc(func(o *options) {
		o.maxOpen = cnt
	})
}

// 初始化参数设置
func initOptions(db *gorm.DB, opts *options) {
	if "" != opts.prefix {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return opts.prefix + defaultTableName
		}
	}

	if true == opts.useLog {
		db.LogMode(opts.useLog)
	}

	if opts.maxLife.Seconds() > 0 {
		db.DB().SetConnMaxLifetime(opts.maxLife)
	}

	if opts.maxIdle > 0 {
		db.DB().SetMaxIdleConns(opts.maxIdle)
	}

	if opts.maxOpen > 0 {
		db.DB().SetMaxOpenConns(opts.maxOpen)
	}
}
