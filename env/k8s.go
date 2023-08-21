//go:build k8s

package env

var Config = config{
	DB: DBConfig{
		DSN: "root:Kezi_520@tcp(webook-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:6380",
	},
}
