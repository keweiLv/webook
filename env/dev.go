//go:build !k8s

package env

var Config = config{
	DB: DBConfig{
		DSN: "localhost:3306",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
