package config

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var v *viper.Viper

var Test bool = os.Getenv("TEST") != ""

func newv() *viper.Viper {
	v := viper.New()
	prefix := ""
	if Test {
		prefix = "../"
	}
	v.SetConfigFile(prefix + "config.ini")
	v.OnConfigChange(func(e fsnotify.Event) {
		loadConfig()
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}

var config = struct {
	host, port           atomic.Pointer[string]
	expiration, maxcount atomic.Int64
	link                 atomic.Pointer[string]
	mode                 atomic.Pointer[string]
	mysqluser            atomic.Pointer[string]
	mysqlpassword        atomic.Pointer[string]
	mysqladdr            atomic.Pointer[string]
}{}

func loadConfig() {
	config.host.Store(ptr(v.GetString("redis.host")))
	config.port.Store(ptr(v.GetString("redis.port")))
	config.expiration.Store(int64(v.GetInt("ip_limit.expiration")))
	config.maxcount.Store(v.GetInt64("ip_limit.maxcount"))
	config.link.Store(ptr(v.GetString("link.path")))
	config.mode.Store(ptr(v.GetString("db.mode")))
	config.mysqluser.Store(ptr(v.GetString("mysql.user")))
	config.mysqlpassword.Store(ptr(v.GetString("mysql.password")))
	config.mysqladdr.Store(ptr(v.GetString("mysql.addr")))
}

func ptr(v string) *string {
	return &v
}

func init() {
	v = newv()
	loadConfig()
}

func GetRedis() (host string, port string) {
	host = *config.host.Load()
	port = *config.port.Load()
	return
}

func GetExpiration() int {
	return int(config.expiration.Load())
}

func GetMaxCount() int64 {
	return config.maxcount.Load()
}

func GetLink() string {
	return *config.link.Load()
}

func GetDbMode() string {
	return *config.mode.Load()
}

func GetDsnInfo() (user, password, addr string) {
	return *config.mysqluser.Load(), *config.mysqlpassword.Load(), *config.mysqladdr.Load()
}
