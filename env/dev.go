//go:build !k8s

package env

var Config = config{
	DB: DBConfig{
		DSN: "root:Kezi_520@tcp(localhost:30001)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:30002",
	},
}
