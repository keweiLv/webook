//go:build k8s

package env

var Config = config{
	DB: DBConfig{
		DSN: "root:Kezi_520@tcp(webook-mysql:6033)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:6379",
	},
}
